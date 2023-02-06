package clickhouse

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	clickhouseV2 "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migrator "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	driverName = "clickhouse"

	msgFailedToAutoMigrate = "failed to automigrate"
	msgFailedToOpenGORM    = "failed to open GORM"
	msgFailedToCloseDB     = "failed to close db"
	msgDBConnected         = "db connected"

	keyDSN = "dsn"

	sourceURLTemplate = "file://%v"
)

var (
	ErrEmptyMigrationsPath = errors.New("migration path should be provided")
)

type (
	Clickhouse struct {
		log  *zerolog.Logger
		sql  *sql.DB
		gorm *gorm.DB
		cfg  Config
	}
)

func New(cfg Config, log zerolog.Logger) *Clickhouse {
	return &Clickhouse{
		log:  &log,
		cfg:  cfg,
		gorm: &gorm.DB{},
	}
}

func (ch *Clickhouse) setupMigrations(ctx context.Context) (err error) {
	var (
		files     []os.FileInfo
		fileBytes []byte
		fPath     string
	)

	if files, err = ioutil.ReadDir(ch.cfg.MigrationsPath); err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fPath = fmt.Sprintf("%v/%v", ch.cfg.MigrationsPath, file.Name())
			if fileBytes, err = ioutil.ReadFile(fPath); err != nil {
				return err
			}

			fileBytes = bytes.ReplaceAll(
				fileBytes,
				[]byte("{{BROKER_SERVER_FOR_KAFKA_ENGINE}}"),
				[]byte(ch.cfg.BrokerServerForKafkaEngine),
			)

			if err = ioutil.WriteFile(fPath, fileBytes, 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ch *Clickhouse) Start(context.Context) (err error) {
	var (
		gormDB *gorm.DB
		sqlDB  *sql.DB
	)

	if err = ch.setupMigrations(context.Background()); err != nil {
		ch.log.Error().Err(err).Msg("failed to setup migrations")
		return err
	}

	sqlDB = clickhouseV2.OpenDB(&clickhouseV2.Options{
		Addr: []string{ch.cfg.Addr},
		Auth: clickhouseV2.Auth{
			Database: ch.cfg.Database,
			Username: ch.cfg.User,
			Password: ch.cfg.Password,
		},
		// Debug: true,
		Settings: clickhouseV2.Settings{
			"max_execution_time": ch.cfg.MaxExecutionTime,
		},
		DialTimeout: ch.cfg.DialTimeout,
	})

	sqlDB.SetMaxIdleConns(ch.cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(ch.cfg.MaxOpenConns)

	gormConfig := &gorm.Config{}
	gormConfig.Logger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	if gormDB, err = gorm.Open(clickhouse.New(clickhouse.Config{Conn: sqlDB}), gormConfig); err != nil {
		ch.log.Error().Err(err).Msg(msgFailedToOpenGORM)
		return err
	}

	ch.sql = sqlDB
	*ch.gorm = *gormDB

	if ch.cfg.AutoMigrate {
		if ch.cfg.MigrationsPath == "" {
			return ErrEmptyMigrationsPath
		}
		err = func() error {
			var (
				driver database.Driver
				m      *migrate.Migrate
			)

			if driver, err = migrator.WithInstance(sqlDB, &migrator.Config{
				MultiStatementEnabled: true,
			}); err != nil {
				return err
			}

			if m, err = migrate.NewWithDatabaseInstance(
				fmt.Sprintf(sourceURLTemplate, ch.cfg.MigrationsPath),
				driverName,
				driver,
			); err != nil {
				return err
			}

			if err = m.Up(); err != nil {
				return err
			}

			return nil
		}()

		if err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				ch.log.Info().Err(err).Msg(msgFailedToAutoMigrate)
			} else {
				ch.log.Error().Err(err).Msg(msgFailedToAutoMigrate)
				return err
			}
		}
	}

	ch.log.Info().Str(keyDSN, ch.cfg.Addr).Msg(msgDBConnected)

	return nil
}

func (ch *Clickhouse) Stop(ctx context.Context) error {
	if err := ch.sql.Close(); err != nil {
		ch.log.Error().Err(err).Msg(msgFailedToCloseDB)
	}

	return nil
}

func (ch *Clickhouse) LatestBlockHeight() (lastHeight int64, err error) {
	if err = ch.gorm.
		Select("height").
		Table("block").
		Order("height DESC").
		Limit(1).
		Scan(&lastHeight).Error; err != nil {
		return 0, err
	}

	return lastHeight, nil
}

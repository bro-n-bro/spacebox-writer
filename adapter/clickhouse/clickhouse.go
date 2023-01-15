package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	clickhouseV2 "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migrator "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	driverName = "clickhouse"
)

type Clickhouse struct {
	log  *zerolog.Logger
	sql  *sql.DB
	gorm *gorm.DB
	cfg  Config
}

func New(cfg Config, log zerolog.Logger) *Clickhouse {
	return &Clickhouse{
		log:  &log,
		cfg:  cfg,
		gorm: &gorm.DB{},
	}
}

func (ch *Clickhouse) Start(context.Context) error {
	sqlDB := clickhouseV2.OpenDB(&clickhouseV2.Options{
		Addr: []string{ch.cfg.Addr},
		Auth: clickhouseV2.Auth{
			Database: ch.cfg.Database,
			Username: ch.cfg.User,
			Password: ch.cfg.Password,
		},
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

	gormDB, err := gorm.Open(clickhouse.New(clickhouse.Config{Conn: sqlDB}), gormConfig)
	if err != nil {
		ch.log.Error().Err(err).Msg("failed to open GORM")
		return err
	}

	ch.sql = sqlDB
	*ch.gorm = *gormDB

	if ch.cfg.AutoMigrate {
		err = func() error {
			var driver database.Driver
			driver, err = migrator.WithInstance(sqlDB, &migrator.Config{
				MultiStatementEnabled: true,
			})
			if err != nil {
				return err
			}
			var m *migrate.Migrate
			m, err = migrate.NewWithDatabaseInstance(
				fmt.Sprintf("file://%v", ch.cfg.MigrationsPath),
				driverName,
				driver,
			)
			if err != nil {
				return err
			}
			err = m.Up()
			if err != nil {
				return err
			}
			return nil
		}()

		if err != nil {
			if err.Error() == "no change" {
				ch.log.Info().Err(err).Msg("failed to automigrate")
			} else {
				ch.log.Error().Err(err).Msg("failed to automigrate")
			}
		}
	}

	ch.log.Info().Str("dsn", ch.cfg.Addr).Msg("db connected")
	return nil
}

func (ch *Clickhouse) Stop(ctx context.Context) error {
	if err := ch.sql.Close(); err != nil {
		ch.log.Error().Err(err).Msg("failed to close db")
	}

	return nil
}

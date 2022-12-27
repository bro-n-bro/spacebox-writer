package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"spacebox-writer/internal/configs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migrator "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Clickhouse struct {
	log  *zerolog.Logger
	sql  *sql.DB
	gorm *gorm.DB
	cfg  configs.Config
}

const (
	driverName = "clickhouse"
	keyCMP     = "cmp"
)

func (clhs *Clickhouse) GetGormDB(ctx context.Context) *gorm.DB { return clhs.gorm }

func New(cfg configs.Config, log zerolog.Logger) *Clickhouse {
	//lg := zerolog.New(os.Stderr).
	//	Output(zerolog.ConsoleWriter{Out: os.Stderr}).
	//	With().
	//	Timestamp().
	//	Str(keyCMP, driverName).Logger()

	return &Clickhouse{
		log:  &log,
		cfg:  cfg,
		gorm: &gorm.DB{},
	}
}

func (clhs *Clickhouse) Start(context.Context) error {
	sqlDB, err := sql.Open(driverName, clhs.cfg.DSN)
	if err != nil {
		clhs.log.Error().Err(err).
			Str("dsn", clhs.cfg.DSN).
			Msg("db connection failed")
		return err
	}

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
		clhs.log.Error().Err(err).Msg("failed to open GORM")
		return err
	}

	clhs.sql = sqlDB
	*clhs.gorm = *gormDB

	if clhs.cfg.AutoMigrate {
		err = func() error {
			driver, err := migrator.WithInstance(sqlDB, &migrator.Config{
				MultiStatementEnabled: true,
			})
			if err != nil {
				return err
			}
			m, err := migrate.NewWithDatabaseInstance(
				fmt.Sprintf("file://%v", clhs.cfg.MigrationsPath),
				fmt.Sprintf(driverName),
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
				clhs.log.Info().Err(err).Msg("failed to automigrate")
			} else {
				clhs.log.Error().Err(err).Msg("failed to automigrate")
			}
		}
	}

	clhs.log.Info().Str("dsn", clhs.cfg.DSN).Msg("db connected")

	return nil
}

func (clhs *Clickhouse) Stop(ctx context.Context) error {
	if err := clhs.sql.Close(); err != nil {
		clhs.log.Error().Err(err).Msg("failed to close db")
	}
	return nil
}

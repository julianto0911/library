package library

import (
	"log"
	"os"
	"testing"
	"time"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func MockGormDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error init mock: %s", err)
		return nil, nil, err
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	gormDB, gerr := gorm.Open(postgres.New(
		postgres.Config{Conn: db}), &gorm.Config{Logger: newLogger})
	if gerr != nil {
		t.Fatalf("error init db: %s", err)
		return nil, nil, err
	}
	gormDB.Logger.LogMode(logger.LogLevel(1))

	return mock, gormDB, nil
}

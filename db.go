package library

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBParam struct {
	Host     string
	Port     string
	Name     string
	Schema   string
	User     string
	Password string
	AppName  string
	Timeout  int
	MaxOpen  int
	MaxIdle  int
	Logging  bool
}

type DBConfiguration struct {
	Host           string
	Port           string
	Schema         string
	DBName         string
	Username       string
	Password       string
	Logging        bool
	SessionName    string
	ConnectTimeOut int
	MaxOpenConn    int
	MaxIdleConn    int
}

func NewDatabaseConnection(config DBConfiguration, l *zap.Logger) (*gorm.DB, error) {
	dbCfg := DBParam{
		Host:     config.Host,
		Port:     config.Port,
		Name:     config.DBName,
		Schema:   config.Schema,
		User:     config.Username,
		Password: config.Password,
		AppName:  config.SessionName,
		Timeout:  config.ConnectTimeOut,
		MaxOpen:  config.MaxOpenConn,
		MaxIdle:  config.MaxIdleConn,
		Logging:  config.Logging,
	}

	level := logger.Silent
	if config.Logging {
		level = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  level,       // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	//sql connection
	sqlConn, err := sql.Open("postgres", makePostgresString(dbCfg))
	if err != nil {
		return nil, errors.Wrap(err, "can't establish db connection")
	}
	sqlConn.SetMaxIdleConns(dbCfg.MaxIdle)
	sqlConn.SetMaxOpenConns(dbCfg.MaxOpen)
	sqlConn.SetConnMaxLifetime(time.Hour)

	db, err := gorm.Open(postgres.New(
		postgres.Config{Conn: sqlConn}),
		&gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, errors.Wrap(err, "can't open db connection")
	}

	return db, err
}

func makePostgresString(p DBParam) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s connect_timeout=%d application_name=%s",
		p.Host, p.Port, p.User, p.Name, p.Password, p.Timeout, p.AppName)
}

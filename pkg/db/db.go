package db

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

type Options struct {
	Type                  Type
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	logger                logger.Interface
}

type Type int8

const (
	MySQL Type = iota
	MSSQL
)

var (
	ErrUnsupportedDBType = errors.New("unsupported db type")
)

func New(options *Options) (db *gorm.DB, err error) {
	zaplogger := zapgorm2.New(zap.L())
	zaplogger.SetAsDefault()
	options.logger = zaplogger

	switch options.Type {
	case MySQL:
		db, err = newMySQL(options)
		break
	case MSSQL:
		db, err = newMSSQL(options)
		break
	default:
		return nil, ErrUnsupportedDBType
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(options.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(options.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(options.MaxIdleConnections)

	return db, nil
}

func newMySQL(options *Options) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		options.Username,
		options.Password,
		options.Host,
		options.Database,
		true,
		"Local")

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: options.logger,
	})

	return db, err
}

func newMSSQL(options *Options) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(`sqlserver://%s:%s@%s?database=%s`,
		options.Username,
		options.Password,
		options.Host,
		options.Database,
	)

	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: options.logger,
	})

	return db, err
}

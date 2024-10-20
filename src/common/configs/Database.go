package configs

import (
	"context"
	"fmt"
	"loyalty-campaigns/src/common/utils"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type IDBConnection interface {
	GetDB() *gorm.DB
	Close() error
}

type dbConnection struct {
	connection *gorm.DB
	logger     utils.ILogger
}

var (
	dbConnInstance *dbConnection
	ddConnOnce     sync.Once
)

func NewDBConnection() IDBConnection {
	ddConnOnce.Do(func() {
		dbConnInstance = &dbConnection{}
		dbConnInstance.logger = utils.NewLogger()
		err := dbConnInstance.connect()
		if err != nil {
			dbConnInstance.logger.Fatal("database connection failed: %v", err)
		}

		err = dbConnInstance.ping()
		if err != nil {
			dbConnInstance.logger.Fatal("failed to ping the database: %v", err)
		}

		dbConnInstance.logger.Success("[OK] database connection")
	})
	return dbConnInstance
}

func (p *dbConnection) connect() error {
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	}

	dsn := "host=localhost user=loyalty_user password=loyalty_pass dbname=loyalty port=5432 sslmode=disable TimeZone=America/Bogota"

	conn, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	p.connection = conn

	return nil
}

func (p *dbConnection) ping() error {
	if p.connection == nil {
		return fmt.Errorf("connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := p.connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func (p *dbConnection) Close() error {
	sqlDB, err := p.connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	err = sqlDB.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

func (p *dbConnection) GetDB() *gorm.DB {
	return p.connection
}

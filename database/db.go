package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxAttempts   int           = 5
	retryWaitTime time.Duration = 10 * time.Second
)

type Parameters struct {
	Host                   string
	Port                   int
	User                   string
	Password               string
	Name                   string
	SkipDefaultTransaction bool
	PrepareStmt            bool
	LogLevel               int
	UseSSL                 bool
}

type Session struct {
	Connection *gorm.DB
	DBName     string
	DB         *sql.DB
}

func NewSession(params *Parameters) (*Session, error) {
	db, err := connectToDB(params)

	if err != nil {
		return nil, fmt.Errorf("could not connect to the database. %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("could not retrieve sql.DB object. %w", err)
	}

	return &Session{
		Connection: db,
		DB:         sqlDB,
		DBName:     params.Name,
	}, nil
}

func connectToDB(params *Parameters) (*gorm.DB, error) {
	sslmode := "disable"

	if params.UseSSL {
		sslmode = "require"
	}

	var db *gorm.DB

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		params.Host, params.User, params.Password, params.Name, params.Port, sslmode)

	dbLogger := logger.New(
		&log.Logger{},
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.LogLevel(params.LogLevel),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	err := retry(maxAttempts, retryWaitTime, func() (err error) {
		// connect to postgres db
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: params.SkipDefaultTransaction,
			PrepareStmt:            params.PrepareStmt,
			Logger:                 dbLogger,
			FullSaveAssociations:   true,
		})

		if err != nil {
			log.Printf("connect to database error: %v", err)
			return err
		}

		// check db connection status
		sqlDB, _ := db.DB()
		err = sqlDB.Ping()

		if err != nil {
			log.Printf("connect to database error: %v", err)
			return err
		}

		return nil
	})

	return db, err
}

func retry(attempts int, sleep time.Duration, f func() error) error {
	var err error

	for i := 0; ; i++ {
		err = f()

		// on success
		if err == nil {
			return nil
		}

		// on no more attempts
		if i >= (attempts - 1) {
			break
		}

		// on fail sleep and retry again
		time.Sleep(sleep)

		log.Println("retrying after error")
	}

	return fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}

// Close databaes session
func (s *Session) Close() error {
	if s == nil {
		return nil
	}

	if err := s.DB.Close(); err != nil {
		return err
	}

	return nil
}

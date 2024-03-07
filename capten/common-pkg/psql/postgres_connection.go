package psql

import (
	"fmt"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	DBMaxOpenConnections    int           = 40
	DBMaxIdleConnections    int           = 10
	DBConnectionMaxLifetime time.Duration = 1 * time.Hour
)

type PostgresDBConn struct {
	DB *gorm.DB
}

// Singleton - Connection poling is handled
var postgresDBConn *PostgresDBConn

var once sync.Once

func getGORMConfig() *gorm.Config {
	return &gorm.Config{
		// Ignoring default transaction started by GORM. Improves performance by upto 30%
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
}

func getPostgresConnString(db string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		"postgres", "jlFlhO8ing", "postgresql.capten.svc.cluster.local", 5432, db)
}

// GetGormClient - Returns db Client
func GetGormClient(connectionStr string) *gorm.DB {
	once.Do(func() {
		db, err := gorm.Open(postgres.Open(connectionStr), getGORMConfig())
		if err != nil {
			panic("Failed to open postgres connection\n" + err.Error())
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic("Failed to open postgres connection\n" + err.Error())
		}

		sqlDB.SetMaxOpenConns(DBMaxOpenConnections)

		sqlDB.SetMaxIdleConns(DBMaxIdleConnections)
		sqlDB.SetConnMaxLifetime(DBConnectionMaxLifetime)
		postgresDBConn = &PostgresDBConn{
			DB: db,
		}
	})

	return postgresDBConn.DB
}

func GetGormWrapper() *gorm.DB {
	return GetGormClient(getPostgresConnString("postgres"))
}

func GetExponentialBackoff() *backoff.ExponentialBackOff {

	expBackoff := backoff.NewExponentialBackOff()

	// MaxElaspedTime - retry mechanism will keep retrying for time set
	expBackoff.MaxElapsedTime = 15 * time.Second

	// MaxIntervalTime - max time between successive retries
	expBackoff.MaxInterval = 4 * time.Second

	// Add random seeding
	expBackoff.Multiplier = 2

	return expBackoff
}

func GetPostgresConnectionStatus() error {
	err := backoff.Retry(func() error {

		client := GetGormClient(getPostgresConnString("postgres"))
		db, err := client.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}, GetExponentialBackoff())
	if err != nil {
		return err
	}
	return nil
}

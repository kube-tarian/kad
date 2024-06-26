package dbinit

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/pkg/errors"
)

type Mode int

const (
	UP    Mode = Mode(1)
	DOWN  Mode = Mode(2)
	PURGE Mode = Mode(3)
)

var log = logging.NewLogger()

type DBConfig struct {
	DBHost     string `envconfig:"PG_DB_HOST" required:"true"`
	DBPort     string `envconfig:"PG_DB_PORT" default:"5432"`
	DBName     string `envconfig:"PG_DB_NAME" required:"true"`
	EntityName string `envconfig:"PG_DB_ENTITY_NAME" default:"postgres"`
	Username   string `envconfig:"PG_DB_SERVICE_USERNAME" required:"true"`
	Password   string `envconfig:"PG_DB_SERVICE_USERPASSWORD" required:"false"`
	SourceURI  string `envconfig:"PG_SOURCE_URI" default:"file:///postgres/migrations"`
}

func RunMigrations(mode Mode) error {
	conf := &DBConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return err
	}

	if len(conf.Password) == 0 {
		serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
			conf.EntityName, conf.Username)
		if err != nil {
			return errors.WithMessage(err, "DB user credential fetching failed")
		}
		conf.Password = serviceCredential.Password
	}
	return RunMigrationsWithConfig(conf, mode)
}

func RunMigrationsWithConfig(conf *DBConfig, mode Mode) error {
	password := url.QueryEscape(conf.Password)
	dbConnectionString := getDbConnectionURLFromDbType(conf, password)
	if err := runMigrations(conf.SourceURI, dbConnectionString, conf.DBName, mode); err != nil {
		return err
	}
	return nil
}

func getDbConnectionURLFromDbType(conf *DBConfig, password string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Username, password, conf.DBHost, conf.DBPort, conf.DBName)
}

func runMigrations(sourceURL, databaseURL, dbName string, mode Mode) (err error) {
	sourceDriver, err := source.Open(sourceURL)
	if err != nil {
		return err
	}
	defer sourceDriver.Close()

	dbDriver, err := database.Open(databaseURL)
	if err != nil {
		return err
	}

	skipMigrations, err := isCurrentSchemaNewerThanLatestMigration(sourceDriver, dbDriver)
	if err != nil {
		return err
	}
	if skipMigrations {
		log.Info("Current DB schema version is newer than the latest available migration. Assuming backward compatible schema and continuing.")
		return nil
	}

	mgr, err := migrate.NewWithInstance("file", sourceDriver, dbName, dbDriver)
	if err != nil {
		return err
	}
	defer mgr.Close()

	switch mode {
	case UP:
		err = mgr.Up()
	case DOWN:
		err = mgr.Down()
	case PURGE:
		err = mgr.Drop()
	}

	if err == migrate.ErrNoChange {
		return nil
	}

	if err != nil {
		return err
	}

	time.Sleep(8 * time.Second)

	return nil
}

func getLatestMigrationVersion(driver source.Driver) (latestVersion uint, err error) {
	currentVersion, err := driver.First()
	if err != nil {
		return
	}

	for {
		latestVersion = currentVersion
		var nextErr error
		currentVersion, nextErr = driver.Next(currentVersion)
		if nextErr != nil {
			if nextErr == os.ErrNotExist {
				return // latest migration reached
			}

			switch nextErr.(type) {
			case *os.PathError:
				return // latest migration reached
			}

			err = nextErr
			return
		}
	}
}

func isCurrentSchemaNewerThanLatestMigration(sourceDriver source.Driver, dbDriver database.Driver) (isNewer bool, err error) {
	dbVersion, _, err := dbDriver.Version()
	if err != nil {
		// assume that the schema_migration table has not been created and continue running the migrations
		log.Info("Unable to get the DB schema version. Continue running the migrations.")
		err = nil
		isNewer = false
		return
	}

	// if no migration has been applied
	if dbVersion == -1 {
		isNewer = false
		return
	}

	latestMigrationVersion, err := getLatestMigrationVersion(sourceDriver)
	if err != nil {
		return
	}

	log.Infof("Current DB version: %d, latest migration version: %d.", dbVersion, latestMigrationVersion)

	isNewer = uint(dbVersion) > latestMigrationVersion
	return
}

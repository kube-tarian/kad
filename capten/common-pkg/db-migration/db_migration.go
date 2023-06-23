// Package migrate contains ...
package dbmigration

import (
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/intelops/go-common/logging"
)

var log = logging.NewLogger()

func RunMigrations(sourceURL, databaseURL, dbName string) (err error) {
	sourceDriver, err := source.Open(sourceURL)
	if err != nil {
		return err
	}

	dbDriver, err := database.Open(databaseURL)
	if err != nil {
		sourceDriver.Close()
		return err
	}

	mgr, err := migrate.NewWithInstance("file", sourceDriver, dbName, dbDriver)
	if err != nil {
		return err
	}
	defer mgr.Close() // this will call close for sourceDriver & dbDriver

	skipMigrations, err := isCurrentSchemaNewerThanLatestMigration(sourceDriver, dbDriver)
	if err != nil {
		return err
	}
	if skipMigrations {
		log.Info("Current DB schema version is newer than the latest available migration. Assuming backward compatible schema and continuing.")
		return nil
	}

	err = mgr.Up()
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

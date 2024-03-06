package migration

import (
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/intelops/go-common/logging"
)

var log = logging.NewLogger()

// This function uses source name "file"
func RunMigrations(sourceURL, databaseURL, dbName string, mode Mode) (err error) {
	sourceDriver, err := source.Open(sourceURL)
	if err != nil {
		return err
	}
	defer sourceDriver.Close()

	return performRunMigrations(sourceDriver, "file", databaseURL, dbName, mode)
}

// This function uses source name "go-bindata"
func RunMigrationsFromBinData(sourcedata *bindata.AssetSource, sourceName, databaseURL, dbName string, mode Mode) (err error) {
	// go-bindata source driver didn't implement Open
	sourceDriver, err := bindata.WithInstance(sourcedata)
	if err != nil {
		return err
	}

	return performRunMigrations(sourceDriver, sourceName, databaseURL, dbName, mode)

}

func performRunMigrations(sourceDriver source.Driver, sourceName, databaseURL, dbName string, mode Mode) error {
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

	mgr, err := getMigrateInstance(sourceName, sourceDriver, dbName, dbDriver, databaseURL)
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

func getMigrateInstance(sourceName string, sourceDriver source.Driver, dbName string, dbDriver database.Driver, databaseURL string) (mgr *migrate.Migrate, err error) {
	switch sourceName {
	case "file":
		mgr, err = migrate.NewWithInstance("file", sourceDriver, dbName, dbDriver)
		if err != nil {
			return nil, err
		}
	case "go-bindata":
		mgr, err = migrate.NewWithSourceInstance("go-bindata", sourceDriver, databaseURL)
		if err != nil {
			return nil, err
		}
	}
	return
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

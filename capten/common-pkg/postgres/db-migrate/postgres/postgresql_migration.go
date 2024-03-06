package postgres

import (
	"context"
	"fmt"
	"net/url"

	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/migration"
	"github.com/pkg/errors"
)

type DBConfig struct {
	DBAddr     string `envconfig:"PG_DB_ADDRESS" required:"true"`
	EntityName string `envconfig:"PG_DB_ENTITY_NAME" required:"true"`
	DBName     string `envconfig:"PG_DB_NAME" required:"true"`
	Username   string `envconfig:"PG_DB_SERVICE_USERNAME" required:"true"`
	Password   string `envconfig:"PG_DB_SERVICE_USER_PASSWORD" required:"false"`
	SourceURI  string `envconfig:"PG_SOURCE_URI" default:"go-bindata"`
}

func RunMigrations(mode migration.Mode) error {
	conf := &DBConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return err
	}

	return RunMigrationsWithConfig(conf, mode)
}

func RunMigrationsWithConfig(conf *DBConfig, mode migration.Mode) error {
	dbConnectionString, err := getDbConnectionURLFromDbType(conf, "")
	if err != nil {
		return errors.WithMessage(err, "DB connection Url create failed")
	}

	// file source based migrations
	if err := migration.RunMigrations(conf.SourceURI, dbConnectionString, conf.DBName, mode); err != nil {
		return err
	}
	return nil
}

func RunMigrationsBinDataWithConfig(conf *DBConfig, s *bindata.AssetSource, mode migration.Mode) error {
	// TODO: password to be passed empty for production
	dbConnectionString, err := getDbConnectionURLFromDbType(conf, conf.Password)
	if err != nil {
		return errors.WithMessage(err, "DB connection Url create failed")
	}

	if err := migration.RunMigrationsFromBinData(s, conf.SourceURI, dbConnectionString, conf.DBName, mode); err != nil {
		return err
	}
	return nil
}

func getDbConnectionURLFromDbType(conf *DBConfig, password string) (string, error) {
	if len(password) == 0 {
		serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
			conf.EntityName, conf.Username)
		if err != nil {
			return "", err
		}
		// postgres://user:password@host:port/dbname?query
		return fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
			conf.Username, url.QueryEscape(serviceCredential.Password), conf.DBAddr, conf.DBName), nil
	}

	// postgres://user:password@host:port/dbname?query
	return fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		conf.Username, password, conf.DBAddr, conf.DBName), nil
}

func GetResource(names []string, afn bindata.AssetFunc) *bindata.AssetSource {
	return bindata.Resource(names, afn)
}

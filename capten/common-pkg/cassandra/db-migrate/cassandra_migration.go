package dbmigrate

import (
	"context"
	"fmt"
	"net/url"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/pkg/errors"
)

type DBConfig struct {
	DBAddrs    []string `envconfig:"DB_ADDRESSES" required:"true"`
	EntityName string   `envconfig:"DB_ENTITY_NAME" required:"true"`
	Keyspace   string   `envconfig:"DB_NAME" required:"true"`
	DBName     string   `envconfig:"DB_NAME" required:"true"` // keyspace
	Username   string   `envconfig:"DB_SERVICE_USERNAME" required:"true"`
	Conistency string   `envconfig:"CONSISTENCY" default:"ALL"`
	SourceURI  string   `envconfig:"SOURCE_URI" default:"file:///cassandra/migrations"`
}

func RunMigrations(log logging.Logger, mode Mode) error {
	conf := &DBConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return err
	}

	dbConnectionString, err := getDbConnectionURLFromDbType(conf)
	if err != nil {
		return errors.WithMessage(err, "DB connection Url create failed")
	}

	if err := runMigrations(conf.SourceURI, dbConnectionString, conf.DBName, mode); err != nil {
		return err
	}
	log.Info("Migrations applied successfully")
	return nil
}

func getDbConnectionURLFromDbType(conf *DBConfig) (dbConnectionURL string, err error) {
	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		conf.EntityName, conf.Username)
	if err != nil {
		return "", err
	}

	passwd := url.QueryEscape(serviceCredential.Password)
	if conf.Conistency == "" {
		err = fmt.Errorf("Cassandra consistency is not provided")
		return
	}
	dbConnectionURL = fmt.Sprintf("cassandra://%s/%s?username=%s&password=%s&consistency=%s",
		conf.DBAddrs[0], conf.DBName, conf.Username, passwd, conf.Conistency)
	return
}

// Package migrate contains ...
package cassandra

import (
	"fmt"
	"net/url"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	dbmigration "github.com/kube-tarian/kad/capten/common-pkg/db-migration"
)

var log = logging.NewLogger()

type DBConfig struct {
	DbDsn       string `envconfig:"DB_ADDRESSES" required:"true" default:"localhost:9042"`
	DbName      string `envconfig:"CASSANDRA_DB_NAME" required:"true"` // keyspace
	Username    string `envconfig:"DB_SERVICE_USERNAME" required:"true" default:"cassandra"`
	Password    string `envconfig:"DB_SERVICE_PASSWD" required:"true" default:"cassandra"`
	Consistency string `envconfig:"CASSANDRA_CONSISTENCY" default:"ALL"`
	SourceURI   string `envconfig:"SOURCE_URI" default:"file:///migrations"`
}

type CassandraMigrate struct {
	conf               *DBConfig
	log                logging.Logger
	sourceURI          string
	dbConnectionString string
}

func NewCassandraMigrate(log logging.Logger) (*CassandraMigrate, error) {
	config := &DBConfig{}
	if err := envconfig.Process("", config); err != nil {
		log.Errorf("Input adminDSN argument or Environment variables are not provided: %v", err)
		return nil, err
	}

	dbConnectionString, err := getDbConnectionURLFromDbType(config)
	if err != nil {
		log.Errorf("Not able form the DB connection Url: %v", err)
		return nil, err
	}

	return &CassandraMigrate{
		conf:               config,
		log:                log,
		sourceURI:          config.SourceURI,
		dbConnectionString: dbConnectionString,
	}, nil
}

func (c *CassandraMigrate) Run(whichDb string, mode dbmigration.Mode) error {
	if err := dbmigration.RunMigrations(c.sourceURI, c.dbConnectionString, whichDb, dbmigration.UP); err != nil {
		log.Errorf("Error string: %s\nError: %+v", err.Error(), err)
		return err
	}
	log.Info("Migrations applied successfully")
	return nil
}

func getDbConnectionURLFromDbType(config *DBConfig) (dbConnectionURL string, err error) {
	passwd := url.QueryEscape(config.Password)
	if config.Consistency == "" {
		err = fmt.Errorf("Cassandra consistency is not provided")
		return
	}
	dbConnectionURL = fmt.Sprintf("cassandra://%s/%s?username=%s&password=%s&consistency=%s",
		config.DbDsn, config.DbName, config.Username, passwd, config.Consistency)

	return
}

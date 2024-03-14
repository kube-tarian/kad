package db

import (
	"context"
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/grpcconn"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/migration"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/postgres"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/source"
)

type DBClient struct {
	conf                *DBConfig
	dbAddress           string
	serviceUserPassword string
}

type DBConfig struct {
	PluginName        string  `envconfig:"PLUGIN_NAME" required:"true"`
	DbOemName         OemName `envconfig:"DB_OEM_NAME" required:"true"`
	DbName            string  `envconfig:"DB_NAME" required:"true"`
	DbServiceUserName string  `envconfig:"DB_SERVICE_USER_NAME" required:"true"`
}

func NewDBClient() (*DBClient, error) {
	conf := &DBConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("DB config read failed, %v", err)
	}

	return NewDBClientWithConfig(conf), nil
}

func NewDBClientWithConfig(conf *DBConfig) *DBClient {
	return &DBClient{
		conf: conf,
	}
}

func (d *DBClient) SetupDatabase() (string, error) {
	client, err := grpcconn.NewGRPCClient()
	if err != nil {
		return "", fmt.Errorf("grpc connection failed while setting up database, %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	response, err := client.SetupDatabase(ctx,
		d.conf.DbOemName.String(),
		d.conf.DbName,
		d.conf.DbServiceUserName,
	)
	if err != nil {
		return "", fmt.Errorf("%s database setup in %s OEM failed: %v", d.conf.DbOemName, d.conf.DbOemName, err)
	}
	fmt.Printf("DB setup status: %v, message: %v, DB URL: %v\n", response.Status, response.StatusMessage, response.VaultPath)

	return response.VaultPath, nil
}

func (d *DBClient) RunMigrations(data map[string][]byte) error {
	binData := source.NewBinData(data)

	err := postgres.RunMigrationsBinDataWithConfig(
		&postgres.DBConfig{
			DBAddr:    d.dbAddress,
			DBName:    d.conf.DbName,
			Username:  d.conf.DbServiceUserName,
			Password:  d.serviceUserPassword,
			SourceURI: "go-bindata",
		},
		postgres.GetResource(binData.FileNames, binData.Asset),
		migration.UP,
	)
	if err != nil {
		return fmt.Errorf("run migrations failed, reason: %v", err)
	}
	return nil
}

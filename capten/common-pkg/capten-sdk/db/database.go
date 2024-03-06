package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/captensdkpb"
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
	DbOemName  OemName `envconfig:"DB_OEM_NAME" required:"true"`
	Parameters *DBParameters
}
type DBParameters struct {
	DbName            string `envconfig:"DB_NAME" required:"true"`
	DbServiceUserName string `envconfig:"DB_SERVICE_USER_NAME" required:"true"`
}

func NewDBClient() (*DBClient, error) {
	oemNameFromEnv, ok := os.LookupEnv("DB_OEM_NAME")
	if !ok {
		return nil, fmt.Errorf("DB_OEM_NAME is missed to provide in environment variables")
	}
	oemName, ok := GetEnum(oemNameFromEnv)
	if !ok {
		return nil, fmt.Errorf("%s: Unsupported database", oemNameFromEnv)
	}

	params := &DBParameters{}
	if err := envconfig.Process(oemName.String(), params); err != nil {
		return nil, fmt.Errorf("DB config parameters read failed, %v", err)
	}

	return NewDBClientWithConfig(&DBConfig{
		DbOemName:  oemName,
		Parameters: params,
	}), nil
}

func NewDBClientWithConfig(conf *DBConfig) *DBClient {
	return &DBClient{
		conf: conf,
	}
}

func (d *DBClient) SetupDatabase() error {
	client, err := grpcconn.NewGRPCClient()
	if err != nil {
		return fmt.Errorf("grpc connection failed while setting up database, %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	response, err := client.Client.SetupDatabase(ctx, &captensdkpb.DBSetupRequest{
		DbOemName:       d.conf.DbOemName.String(),
		DbName:          d.conf.Parameters.DbName,
		ServiceUserName: d.conf.Parameters.DbServiceUserName,
	})
	if err != nil {
		return fmt.Errorf("%s database setup in %s OEM failed: %v", d.conf.DbOemName, d.conf.DbOemName, err)
	}
	fmt.Printf("DB setup status: %v, message: %v, DB URL: %v\n", response.Status, response.StatusMessage, response.DbURL)

	return nil
}

func (d *DBClient) RunMigrations(data map[string][]byte) error {
	binData := source.NewBinData(data)

	err := postgres.RunMigrationsBinDataWithConfig(
		&postgres.DBConfig{
			DBAddr:    d.dbAddress,
			DBName:    d.conf.Parameters.DbName,
			Username:  d.conf.Parameters.DbServiceUserName,
			Password:  d.serviceUserPassword,
			SourceURI: "go-bindata",
		},
		postgres.GetResource(binData.FileNames, binData.Asset),
		migration.UP,
	)
	if err != nil {
		return fmt.Errorf("Run migrations failed, reason: %v", err)
	}
	return nil
}

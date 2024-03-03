package db

import (
	"context"
	"fmt"
	"os"
	"postgressetup/pkg/common-pkg/sdk-capten/grpcconn"
	"postgressetup/pkg/db-migrate/source"
	"postgressetup/server/pkg/pb/serverpb"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type SdkDbClient interface {
}

type DBClient struct {
	conf *DBConfig
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

	response, err := client.Client.SetupDatabase(ctx, &serverpb.DBSetupRequest{
		DbOemName:       d.conf.DbOemName.String(),
		DbName:          d.conf.Parameters.DbName,
		ServiceUserName: d.conf.Parameters.DbServiceUserName,
	})
	if err != nil {
		return fmt.Errorf("%s database setup in %s OEM failed: %v", d.conf.DbOemName, d.conf.DbOemName, err)
	}
	fmt.Printf("DB setup success: %v", response)

	return nil
}

func (d *DBClient) RunMigrations(migrationsPath string) error {
	client, err := grpcconn.NewGRPCClient()
	if err != nil {
		return fmt.Errorf("grpc connection failed while applying migrations, %v", err)
	}
	defer client.Close()

	var migrations = map[string][]byte{}
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("reading migration scripts from path %s failed: %v", migrationsPath, err)
	}

	for _, file := range files {
		migrations[file.Name()], err = getFileContent(migrationsPath, file.Name())
		if err != nil {
			return fmt.Errorf("reading the migrations file %s content failed, %v", file.Name(), err)
		}
	}

	binData := source.NewBinData(migrations)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	response, err := client.Client.RunMigrations(ctx, &serverpb.DBMigrationRequest{
		DbOemName:       d.conf.DbOemName.String(),
		DbName:          d.conf.Parameters.DbName,
		ServiceUserName: d.conf.Parameters.DbServiceUserName,
		Migrations:      binData.FilesMap,
	})
	if err != nil {
		return fmt.Errorf("%s migrations applying in %s OEM failed: %v", d.conf.DbOemName, d.conf.DbOemName, err)
	}
	fmt.Printf("migrations applied successfully: %v", response)
	return nil
}

func getFileContent(migrationsPath, fileName string) ([]byte, error) {
	return os.ReadFile(migrationsPath + fileName)
}

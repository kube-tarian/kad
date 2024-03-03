package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/server/pkg/pb/serverpb"

	"github.com/kube-tarian/kad/capten-sdk/db"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/migration"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/postgres"
	"github.com/kube-tarian/kad/capten/common-pkg/postgres/db-migrate/source"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
)

func (s *Server) SetupDatabase(ctx context.Context, req *serverpb.DBSetupRequest) (*serverpb.DBSetupResponse, error) {
	conf, err := readConfig()
	if err != nil {
		s.log.Error(err.Error())
		return &serverpb.DBSetupResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	s.log.Info("Creating new db for configuration")
	// Setup the database in postgres
	switch req.DbOemName {
	case db.POSTGRES.String():
		err = setupPostgresDatabase(s.log, conf)
		if err != nil {
			s.log.Error(err.Error())
			return &serverpb.DBSetupResponse{
				Status:        serverpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: err.Error(),
			}, nil
		}
	default:
		s.log.Error("Unsupported Database OEM %s", req.DbOemName)
		return &serverpb.DBSetupResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: fmt.Sprintf("Unsupported Database OEM %s", req.DbOemName),
		}, nil
	}
	s.log.Infof("Setup of new db %s is Done", conf.DBName)

	return &serverpb.DBSetupResponse{
		Status:              serverpb.StatusCode_OK,
		StatusMessage:       "Database setup in postgres succesful",
		ServiceUserPassword: conf.Password,
		DbURL:               conf.DBAddress + ":5432",
	}, nil
}

func (s *Server) RunMigrations(ctx context.Context, req *serverpb.DBMigrationRequest) (*serverpb.DBMigrationResponse, error) {
	conf, err := readConfig()
	if err != nil {
		return &serverpb.DBMigrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	// Run migrations
	switch req.DbOemName {
	case db.POSTGRES.String():
		err = runMigrations(s.log, conf, req)
		if err != nil {
			s.log.Error(err.Error())
			return &serverpb.DBMigrationResponse{
				Status:        serverpb.StatusCode_INTERNRAL_ERROR,
				StatusMessage: err.Error(),
			}, nil
		}
	default:
		s.log.Error("Unsupported Database OEM %s", req.DbOemName)
		return &serverpb.DBMigrationResponse{
			Status:        serverpb.StatusCode_INTERNRAL_ERROR,
			StatusMessage: fmt.Sprintf("Unsupported Database OEM %s", req.DbOemName),
		}, nil
	}

	return &serverpb.DBMigrationResponse{
		Status:        serverpb.StatusCode_OK,
		StatusMessage: "run migrations sucessful",
	}, nil
}

// Read the Postgres DB configuration
func readConfig() (*dbinit.Config, error) {
	conf := &dbinit.Config{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func setupPostgresDatabase(log logging.Logger, conf *dbinit.Config) error {
	return dbinit.CreatedDatabaseWithConfig(log, conf)
}

func runMigrations(log logging.Logger, conf *dbinit.Config, req *serverpb.DBMigrationRequest) error {
	binData := source.NewBinData(req.Migrations)

	err := postgres.RunMigrationsBinDataWithConfig(
		log,
		&postgres.DBConfig{
			DBAddr:     conf.DBAddress,
			EntityName: conf.EntityName,
			DBName:     conf.DBName,
			Username:   conf.DBServiceUsername,
			Password:   conf.Password,
			SourceURI:  "go-bindata",
		},
		postgres.GetResource(binData.FileNames, binData.Asset),
		migration.UP,
	)
	if err != nil {
		log.Errorf("Run migrations failed, reason: %v", err)
		return err
	}
	return nil
}

package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captensdkpb"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
)

func (a *Agent) SetupDatabase(ctx context.Context, req *captensdkpb.DBSetupRequest) (*captensdkpb.DBSetupResponse, error) {
	conf, err := readConfig()
	if err != nil {
		a.log.Error(err.Error())
		return &captensdkpb.DBSetupResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}

	a.log.Info("Creating new db for configuration")
	// Setup the database in postgres
	switch req.DbOemName {
	// case db.POSTGRES.String():
	case "POSTGRES":
		err = setupPostgresDatabase(a.log, conf)
		if err != nil {
			a.log.Error(err.Error())
			return &captensdkpb.DBSetupResponse{
				Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
				StatusMessage: err.Error(),
			}, nil
		}
	default:
		a.log.Error("Unsupported Database OEM %s", req.DbOemName)
		return &captensdkpb.DBSetupResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: fmt.Sprintf("Unsupported Database OEM %s", req.DbOemName),
		}, nil
	}
	a.log.Infof("Setup of new db %s is Done", conf.DBName)

	return &captensdkpb.DBSetupResponse{
		Status:              captensdkpb.StatusCode_OK,
		StatusMessage:       "Database setup in postgres succesful",
		ServiceUserPassword: conf.Password,
		DbURL:               conf.DBAddress + ":5432",
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

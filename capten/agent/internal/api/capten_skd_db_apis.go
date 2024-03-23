package api

import (
	"context"
	"fmt"

	"github.com/kube-tarian/kad/capten/agent/internal/pb/captensdkpb"
	"github.com/kube-tarian/kad/capten/common-pkg/capten-sdk/db"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
)

func (a *Agent) SetupDatabase(ctx context.Context, req *captensdkpb.SetupDatabaseRequest) (*captensdkpb.SetupDatabaseResponse, error) {
	a.log.Info("Creating new db for configuration")
	var vaultPath string
	var err error
	// Setup the database in postgres
	switch req.DbOemName {
	case db.POSTGRES.String():
		vaultPath, err = setupPostgresDatabase(a.log, req)
	default:
		err = fmt.Errorf("unsupported Database OEM %s", req.DbOemName)
	}
	if err != nil {
		a.log.Error(err.Error())
		return &captensdkpb.SetupDatabaseResponse{
			Status:        captensdkpb.StatusCode_INTERNAL_ERROR,
			StatusMessage: err.Error(),
		}, nil
	}
	a.log.Infof("Setup of new db %s is Done", req.DbName)

	return &captensdkpb.SetupDatabaseResponse{
		Status:                        captensdkpb.StatusCode_OK,
		StatusMessage:                 "Database setup in postgres succesful",
		SvcUserCredentialsPathInVault: vaultPath,
	}, nil
}

// Read the Postgres DB configuration
func readConfig() (*dbinit.Config, error) {
	var baseConfig dbinit.BaseConfig
	if err := envconfig.Process("", &baseConfig); err != nil {
		return nil, err
	}
	return &dbinit.Config{
		BaseConfig: baseConfig,
	}, nil
}

func setupPostgresDatabase(log logging.Logger, req *captensdkpb.SetupDatabaseRequest) (vaultPath string, err error) {
	conf, err := readConfig()
	if err != nil {
		log.Error(err.Error())
		return
	}

	conf.DBName = req.DbName
	conf.DBServiceUsername = req.ServiceUserName
	conf.Password = dbinit.GenerateRandomPassword(12)

	err = dbinit.CreatedDatabaseWithConfig(log, conf)
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Insert into vault path plugin/<plugin-name>/<svc-entity> => plugin/test/postgres
	cred := credentials.PrepareServiceCredentialMap(credentials.ServiceCredential{
		UserName: conf.DBServiceUsername,
		Password: conf.Password,
		AdditionalData: map[string]string{
			"db-url":  conf.DBAddress,
			"db-name": conf.DBName,
		},
	})
	return fmt.Sprintf("%s/%s/%s", credentials.CertCredentialType, req.PluginName, conf.EntityName),
		credential.PutPluginCredential(context.TODO(), req.PluginName, conf.EntityName, cred)
}

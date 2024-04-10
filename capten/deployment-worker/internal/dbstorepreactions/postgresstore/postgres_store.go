package postgresstore

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
	dbinit "github.com/kube-tarian/kad/capten/common-pkg/postgres/db-init"
	"github.com/kube-tarian/kad/capten/deployment-worker/internal/k8sops"
)

const (
	dbNameTemplate = "%s-db"
)

func SetupPostgresDatabase(log logging.Logger, pluginName, namespace, cmName string, k8sClient *k8s.K8SClient) error {
	conf, err := readConfig()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	vaultPath := fmt.Sprintf("%s/%s/%s", credentials.CertCredentialType, pluginName, conf.EntityName)

	conf.DBName = fmt.Sprintf(dbNameTemplate, pluginName)
	conf.DBServiceUsername = pluginName
	conf.Password = dbinit.GenerateRandomPassword(12)

	err = dbinit.CreatedDatabaseWithConfig(log, conf)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	// Insert into vault path plugin/<plugin-name>/<svc-entity> => plugin/test/postgres
	cred := credentials.PrepareServiceCredentialMap(credentials.ServiceCredential{
		UserName: conf.DBServiceUsername,
		Password: conf.Password,
		AdditionalData: map[string]string{
			"db-url":       conf.DBAddress,
			"db-name":      conf.DBName,
			"service-user": pluginName,
		},
	})

	k8sops.CreateUpdateConfigmap(context.TODO(), log, namespace, cmName, map[string]string{
		"vault-path": vaultPath,
	}, k8sClient)
	if err != nil {
		return err
	}

	return credential.PutPluginCredential(context.TODO(), pluginName, conf.EntityName, cred)
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

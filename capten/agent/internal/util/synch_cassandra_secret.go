package util

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

type SyncSecretConfig struct {
	DBAdminCredIdentifier string `envconfig:"DB_ADMIN_CRED_IDENTIFIER" default:"cassandra-admin"`
	EntityName            string `envconfig:"DB_ENTITY_NAME" required:"true"`
	SecretName            string `envconfig:"CASSANDRA_SECRET_NAME" required:"true"`
	Namespace             string `envconfig:"POD_NAMESPACE" required:"true"`
}

func SyncCassandraAdminSecret(log logging.Logger) error {
	conf := &SyncSecretConfig{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("cassandra config read failed, %v", err)
	}

	k8sClient, err := k8s.NewK8SClient(log)
	if err != nil {
		return err
	}

	res, err := k8sClient.GetSecretData(conf.Namespace, conf.SecretName)
	if err != nil {
		return err
	}

	userName := res.Data["username"]
	password := res.Data["password"]
	if len(userName) == 0 || len(password) == 0 {
		return fmt.Errorf("credentials not found in the secret")
	}

	err = credential.PutServiceUserCredential(context.Background(), conf.EntityName, conf.DBAdminCredIdentifier, userName, password)
	return err
}

package dbinit

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
	"github.com/pkg/errors"
)

const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	numberChars    = "0123456789"
	specialChars   = "!@#$%^&"
)

type Config struct {
	DBAddresses           []string `envconfig:"DB_ADDRESSES" required:"true"`
	DBAdminCredIdentifier string   `envconfig:"DB_ADMIN_CRED_IDENTIFIER" default:"cassandra-admin"`
	DBReplicationFactor   string   `envconfig:"DB_REPLICATION_FACTOR" required:"true" default:"1"`
	EntityName            string   `envconfig:"DB_ENTITY_NAME" required:"true"`
	DBName                string   `envconfig:"DB_NAME" required:"true"`
	DBServiceUsername     string   `envconfig:"DB_SERVICE_USERNAME" required:"true"`
}

func CreatedDatabase(log logging.Logger) (err error) {
	log.Debug("Creating new db for configuration")
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("cassandra config read faile, %v", err)
	}
	if len(conf.DBAddresses) == 0 {
		return errors.New("cassandra DB addresses are empty")
	}

	adminCredential, err := credential.GetServiceUserCredential(context.Background(),
		conf.EntityName, conf.DBAdminCredIdentifier)
	if err != nil {
		return err
	}

	adminClient, err := NewCassandraAdmin(log, conf.DBAddresses, adminCredential.UserName, adminCredential.Password)
	if err != nil {
		return
	}
	defer adminClient.Close()

	log.Info("Creating lock schema change db")
	err = adminClient.CreateLockSchemaDb(conf.DBReplicationFactor)
	if err != nil {
		return
	}

	log.Infof("Creating new db %s with %s", conf.DBName, conf.DBReplicationFactor)
	err = adminClient.CreateDb(conf.DBName, conf.DBReplicationFactor)
	if err != nil {
		return
	}

	var serviceUserPassword string
	var serviceUserName string
	serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
		conf.EntityName, conf.DBServiceUsername)
	if err != nil {
		log.Infof("user %s not exist in DB, %v", conf.DBServiceUsername, err)
		serviceUserPassword = generateRandomPassword(12)
		serviceUserName = conf.DBServiceUsername
	} else {
		serviceUserPassword = serviceCredential.Password
		serviceUserName = serviceCredential.UserName
	}

	log.Infof("Creating new service user %s", serviceUserName)
	err = adminClient.CreateDbUser(serviceUserName, serviceUserPassword)
	if err != nil {
		return
	}

	err = credential.PutServiceUserCredential(context.Background(), conf.EntityName,
		conf.DBServiceUsername, conf.DBServiceUsername, serviceUserPassword)
	if err != nil {
		return
	}

	log.Info("Grant permission to service user")
	err = adminClient.GrantPermission(conf.DBServiceUsername, conf.DBName)
	if err != nil {
		return
	}
	return
}

func generateRandomPassword(length int) string {
	var passwordChars = uppercaseChars + lowercaseChars + numberChars + specialChars
	password := make([]byte, length)
	maxCharIndex := big.NewInt(int64(len(passwordChars)))

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, maxCharIndex)
		password[i] = passwordChars[randomIndex.Int64()]
	}

	return string(password)
}

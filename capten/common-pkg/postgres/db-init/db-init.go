package dbinit

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/intelops/go-common/credentials"
	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/credential"
)

const (
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	numberChars    = "0123456789"
	specialChars   = "!@#$%^&"
)

type Config struct {
	BaseConfig
	DBName            string `envconfig:"PG_DB_NAME" required:"false"`
	DBServiceUsername string `envconfig:"PG_DB_SERVICE_USERNAME" required:"false"`
	Password          string `envconfig:"PG_DB_SERVICE_USER_PASSWORD" required:"false"`
	AdminPassword     string `envconfig:"PG_DB_ADMIN_PASSWORD" required:"false"`
}

type BaseConfig struct {
	DBAddress             string `envconfig:"PG_DB_ADDRESS" required:"true"`
	DBAdminCredIdentifier string `envconfig:"PG_DB_ADMIN_CRED_IDENTIFIER" default:"postgres-admin"`
	EntityName            string `envconfig:"PG_DB_ENTITY_NAME" required:"true"`
}

func CreatedDatabase(log logging.Logger) (err error) {
	log.Debug("Creating new db for configuration")
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return fmt.Errorf("postgres config read faile, %v", err)
	}

	return CreatedDatabaseWithConfig(log, conf)
}

func CreatedDatabaseWithConfig(log logging.Logger, conf *Config) (err error) {
	var adminCredential credentials.ServiceCredential
	if len(conf.AdminPassword) == 0 {
		adminCredential, err = credential.GetServiceUserCredential(context.Background(), conf.EntityName, conf.DBAdminCredIdentifier)
		if err != nil {
			return err
		}
	} else {
		adminCredential = credentials.ServiceCredential{
			UserName: "postgres",
			Password: conf.AdminPassword,
		}
	}

	adminClient, err := NewPostgresAdmin(log, conf.DBAddress, adminCredential.UserName, adminCredential.Password)
	if err != nil {
		return
	}

	err = adminClient.CreateDb(conf.DBName)
	if err != nil {
		return
	}

	var serviceUserPassword string
	var serviceUserName string
	if len(conf.Password) == 0 {
		serviceCredential, err := credential.GetServiceUserCredential(context.Background(),
			conf.EntityName, conf.DBServiceUsername)
		if err != nil {
			log.Infof("user %s not exist in DB, %v", conf.DBServiceUsername, err)
			serviceUserPassword = GenerateRandomPassword(12)
			serviceUserName = conf.DBServiceUsername
		} else {
			serviceUserPassword = serviceCredential.Password
			serviceUserName = serviceCredential.UserName
		}
	} else {
		serviceUserPassword = conf.Password
		serviceUserName = conf.DBServiceUsername
	}

	log.Infof("Creating new service user %s", serviceUserName)
	err = adminClient.CreateDbUser(serviceUserName, serviceUserPassword)
	if err != nil {
		return
	}

	if len(conf.Password) == 0 {
		err = credential.PutServiceUserCredential(context.Background(), conf.EntityName,
			conf.DBServiceUsername, conf.DBServiceUsername, serviceUserPassword)
		if err != nil {
			return
		}
	}

	log.Info("Grant permission to service user")
	err = adminClient.GrantPermission(conf.DBServiceUsername, conf.DBAddress, adminCredential.UserName, adminCredential.Password, conf.DBName)
	if err != nil {
		return
	}
	return
}

func GenerateRandomPassword(length int) string {
	var passwordChars = uppercaseChars + lowercaseChars + numberChars + specialChars
	password := make([]byte, length)
	maxCharIndex := big.NewInt(int64(len(passwordChars)))

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, maxCharIndex)
		password[i] = passwordChars[randomIndex.Int64()]
	}

	return string(password)
}

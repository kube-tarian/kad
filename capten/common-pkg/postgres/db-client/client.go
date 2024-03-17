package dbclient

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	DSNTemplate = "user=%s password=%s dbname=%s host=%s port=%s sslmode=%s"
)

type Config struct {
	DBAddr              string `envconfig:"DB_ADDRESS" required:"true"`
	DBPort              string `envconfig:"DB_Port" required:"true"`
	ServiceUsername     string `envconfig:"DB_SERVICE_USERNAME" required:"true"`
	ServiceUserPassword string `envconfig:"DB_SERVICE_PASSWORD" required:"false"`
	EntityName          string `envconfig:"DB_ENTITY_NAME" required:"true"`
	DBName              string `envconfig:"DB_NAME" required:"true"`
	MaxRetryCount       int    `envconfig:"MAX_RETRY_COUNT" default:"3"`
	MaxConnectionCount  int    `envconfig:"MAX_CLUSTER_CONNECTION_COUNT" default:"5"`
	SSLMode             string `envconfig:"SSL_MODE" default:"disable"`
}

type Client struct {
	Conf *Config
	db   *gorm.DB
}

func NewClient() (*Client, error) {
	conf := &Config{}
	if err := envconfig.Process("", conf); err != nil {
		return nil, fmt.Errorf("cassandra config read faile, %v", err)
	}

	dsn := fmt.Sprintf(DSNTemplate, conf.ServiceUsername, conf.ServiceUserPassword, conf.DBName, conf.DBAddr, conf.DBPort, conf.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Client{
		Conf: conf,
		db:   db,
	}, nil
}

func (c *Client) GetDB() *gorm.DB {
	return c.db
}

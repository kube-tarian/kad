package config

import (
	"github.com/kelseyhightower/envconfig"
)

type SericeConfig struct {
	Host                      string `envconfig:"HOST" default:"0.0.0.0"`
	Port                      int    `envconfig:"PORT" default:"9091"`
	RestPort                  int    `envconfig:"REST_PORT" default:"8443"`
	Mode                      string `envconfig:"MODE" default:"production"`
	AuthEnabled               bool   `envconfig:"AUTH_ENABLED" default:"false"`
	CrossplaneSyncJobEnabled  bool   `envconfig:"CROSSPLANE_SYNC_JOB_ENABLED" default:"true"`
	CrossplaneSyncJobInterval string `envconfig:"CROSSPLANE_SYNC_JOB_INTERVAL" default:"@every 5m"`
	TektonSyncJobEnabled      bool   `envconfig:"TEKTON_SYNC_JOB_ENABLED" default:"true"`
	TektonSyncJobInterval     string `envconfig:"TEKTON_SYNC_JOB_INTERVAL" default:"@every 1h"`
	DomainName                string `envconfig:"DOMAIN_NAME" default:"example.com"`
	ClusterCAIssuerName       string `envconfig:"AGENT_CLUSTER_CA_ISSUER_NAME" default:"agent-ca-issuer"`
}

func GetServiceConfig() (*SericeConfig, error) {
	cfg := &SericeConfig{}
	err := envconfig.Process("", cfg)
	return cfg, err
}

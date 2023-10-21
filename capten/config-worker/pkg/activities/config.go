package activities

import "github.com/kelseyhightower/envconfig"

type Config struct {
	GitDefaultCommiterName  string `envconfig:"GIT_COMMIT_NAME" default:"capten-bot"`
	GitDefaultCommiterEmail string `envconfig:"GIT_COMMIT_EMAIL" default:"capten-bot@intelops.dev"`
	VaultEntityName         string `envconfig:"VAULT_ENTITY_NAME" default:"gitproject"`
	GitCLoneDir             string `envconfig:"GIT_CLONE_DIR" default:"/gitCloneDir"`
	CiCDTemplateRepo        string `envconfig:"CI_CD_TEMPLATE_REPO" default:"https://github.com/intelops/capten-templates.git"`
}

func GetConfig() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

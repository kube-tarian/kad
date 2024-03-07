package postgres

import (
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/psql"
	"gorm.io/gorm"
)

type Postgres struct {
	db  *gorm.DB
	log logging.Logger
}

func NewPostgres(log logging.Logger) *Postgres {
	db := psql.GetGormWrapper()
	return &Postgres{log: log, db: db}
}

func (handler *Postgres) RunMigrations() error {

	err := handler.db.AutoMigrate(GitProjects{}, CloudProviders{}, ContainerRegistry{})

	return err
}

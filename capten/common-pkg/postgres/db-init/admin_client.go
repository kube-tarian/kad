package dbinit

import (
	"fmt"
	"strings"

	"github.com/intelops/go-common/logging"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	createDatabaseQuery          = `CREATE DATABASE %s`
	createUserQuery              = "CREATE USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	alterUserQuery               = "ALTER USER %s WITH PASSWORD '%s' NOSUPERUSER;"
	grantPermissionDatabaseQuery = "GRANT ALL privileges ON database %s TO %s ;"
	grantSchemaPublicQuery       = "GRANT ALL ON SCHEMA public TO %s;"
	dsnTemplate                  = "postgres://%s:%s@%s/"
	dsnDBTemplate                = "postgres://%s:%s@%s/%s"
)

type PostgresAdmin struct {
	log     logging.Logger
	session *gorm.DB
}

func NewPostgresAdmin(logger logging.Logger, dbAddr, dbAdminUsername, dbAdminPassword string) (*PostgresAdmin, error) {
	pg := &PostgresAdmin{
		log: logger,
	}
	err := pg.initSession(dbAddr, dbAdminUsername, dbAdminPassword)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

func (p *PostgresAdmin) initSession(dbAddr, dbAdminUsername, dbAdminPassword string) (err error) {
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable"
	dsn := fmt.Sprintf(dsnTemplate, dbAdminUsername, dbAdminPassword, dbAddr)
	p.session, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return
}

func (p *PostgresAdmin) Session() *gorm.DB {
	return p.session
}

func (p *PostgresAdmin) CreateDbUser(serviceUsername string, servicePassword string) (err error) {
	result := p.session.Exec(fmt.Sprintf(createUserQuery, serviceUsername, servicePassword))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "already exists") {
			return p.updateDbUser(serviceUsername, servicePassword)
		} else {
			err = errors.WithMessage(err, "failed to create service user")
			return
		}
	}
	return
}

func (p *PostgresAdmin) GrantPermission(serviceUsername string, dbAddr, dbAdminUsername, dbAdminPassword, dbName string) error {
	dbDSN := fmt.Sprintf(dsnDBTemplate, dbAdminUsername, dbAdminPassword, dbAddr, dbName)
	session, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		return errors.WithMessage(err, "failed to grant permission to service user")
	}

	result := session.Exec(fmt.Sprintf(grantPermissionDatabaseQuery, dbName, serviceUsername))
	if result.Error != nil {
		return errors.WithMessage(result.Error, "failed to grant permission to service user")
	}

	result = session.Exec(fmt.Sprintf(grantSchemaPublicQuery, serviceUsername))
	if result.Error != nil {
		return errors.WithMessage(result.Error, "failed to grant permission to service user")
	}
	p.log.Infof("Grant permission successful")
	return nil
}

func (p *PostgresAdmin) CreateDb(dbName string) (err error) {
	result := p.session.Exec(fmt.Sprintf(createDatabaseQuery, dbName))
	if result.Error != nil {
		if !strings.Contains(result.Error.Error(), "already exists") {
			err = errors.WithMessage(err, "failed to create the Database")
		}
	}
	return
}

func (c *PostgresAdmin) updateDbUser(serviceUsername string, servicePassword string) (err error) {
	result := c.session.Exec(fmt.Sprintf(alterUserQuery, serviceUsername, servicePassword))
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "already exists") {
			return c.updateDbUser(serviceUsername, servicePassword)
		} else {
			err = errors.WithMessage(err, "failed to create service user")
			return
		}
	}
	return
}

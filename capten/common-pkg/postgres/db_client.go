// Package postgresdb ...
package postgresdb

import (
	"errors"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ObjectNotExist       gerrors.ErrorCode = "Object Not Exists"
	DuplicateRecord      gerrors.ErrorCode = "Record already Exist"
	AuthenticationFailed gerrors.ErrorCode = "Authentication Failed"
	PostgresDBError      gerrors.ErrorCode = "Postgres DB Error"
	NoError              gerrors.ErrorCode = "No Error"
)

const (
	DBConnectError      gerrors.ErrorCode = "Error connecting to the DB"
	GormGetDialectError gerrors.ErrorCode = "Error connecting to the DB"
)

type Config struct {
	Username     string `envconfig:"PG_DB_SERVICE_USERNAME" required:"true"`
	Password     string `envconfig:"PG_DB_SERVICE_USERPASSWORD" required:"true"`
	DBHost       string `envconfig:"PG_DB_HOST" required:"true"`
	DBPort       string `envconfig:"PG_DB_PORT" required:"true"`
	DatabaseName string `envconfig:"PG_DB_NAME" required:"true"`
	IsTLSEnabled bool   `envconfig:"PG_DB_TLS_ENABLED" default:"false"`
}

type DBClient struct {
	session *gorm.DB // this will be initialized later as part of lazy
	logger  logging.Logger
}

func NewDBFromENV(logger logging.Logger) (*gorm.DB, error) {
	conf := Config{}
	if err := envconfig.Process("", &conf); err != nil {
		return nil, err
	}

	return NewDB(&conf, logger)
}

func NewDB(conf *Config, logger logging.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(DataSourceName(conf.Username, conf.Password, fmt.Sprintf("%s:%s", conf.DBHost, conf.DBPort), conf.DatabaseName, conf.IsTLSEnabled)), &gorm.Config{})
	if err != nil {
		return nil, gerrors.NewFromError(DBConnectError, err)
	}
	// db.Session(&gorm.Session{Logger: logger})
	return db, nil
}

func DataSourceName(username string, password string, address string, dbName string, isTLSEnabled bool) string {
	return fmt.Sprintf("postgres://%s:%s@%s/%s",
		username, password, address, dbName)
}

// NewDBClient returns new DB client instance
func NewDBClient(logger logging.Logger) (store *DBClient, err error) {
	logger.Debug("Getting db connection for ...")
	session, err := NewDBFromENV(logger)
	if err != nil {
		return nil, fmt.Errorf("error while creating mariadb client session, %v", err)
	}
	store = &DBClient{
		session: session,
		logger:  logger,
	}
	return store, nil
}

func (db *DBClient) Session() *gorm.DB {
	return db.session
}

func (db *DBClient) getErrorCode(err error) gerrors.ErrorCode {
	if err == nil {
		return NoError
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gerrors.NotFound
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return DuplicateRecord
	}

	return PostgresDBError
}

// Create insert the value into database
func (db *DBClient) Create(value interface{}) (err error) {
	err = db.session.Create(value).Error
	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// Save update value in database, if the value doesn't have primary key, will insert it
func (db *DBClient) Save(value interface{}) (err error) {
	err = db.session.Save(value).Error
	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// Delete delete value match given conditions, if the value has primary key, then will including the primary key as condition
// Save update value in database, if the value doesn't have primary key, will insert it
func (db *DBClient) Delete(value interface{}, where interface{}) (err error) {
	result := db.session.Where(where).Delete(value)
	if result.Error != nil {
		err = result.Error
		return
	}
	if result.RowsAffected == 0 && result.Error == nil {
		err = gerrors.New(ObjectNotExist, "")
	}
	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

func (db *DBClient) FindFirst(value interface{}, where interface{}, args ...interface{}) (err error) {
	err = db.session.Where(where, args).First(value).Error
	if err == nil {
		return
	}

	if gorm.ErrRecordNotFound == err {
		err = gerrors.New(ObjectNotExist, "")
	}

	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// Find find records that match given conditions
func (db *DBClient) Find(value interface{}, where interface{}, args ...interface{}) (err error) {
	if where != nil {
		err = db.session.Where(where, args).Find(value).Error
	} else {
		err = db.session.Find(value).Error
	}
	if err == nil {
		return
	}

	if gorm.ErrRecordNotFound == err {
		err = gerrors.New(ObjectNotExist, "")
	}

	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// FindWithOrder Find find records that match given conditions with the provided order
func (db *DBClient) FindWithOrder(value interface{}, order string, where interface{}) (err error) {
	if where != nil {
		err = db.session.Where(where).Order(order).Find(value).Error
	} else {
		err = db.session.Order(order).Find(value).Error
	}
	if err == nil {
		return
	}
	if gorm.ErrRecordNotFound == err {
		err = gerrors.New(ObjectNotExist, "")
	}
	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// Update the value into database
func (db *DBClient) Update(value interface{}, where interface{}) (err error) {
	// First do the DB update operation
	err = db.session.Model(value).Where(where).Updates(value).Error
	if err == nil {
		return
	}
	if gorm.ErrRecordNotFound == err {
		err = gerrors.New(ObjectNotExist, "")
	}
	if err != nil {
		err = db.checkError(err)
		return
	}
	return
}

// Clear the Rule Table
func (db *DBClient) Clear() {
	// clearing rule data in store
	err := db.session.Exec("TRUNCATE TABLE rule").Error
	if err != nil {
		db.logger.Errorf("Clearing rule table failed with error %v", err)
	}
}

func (db *DBClient) checkError(operation error) (err error) {
	err = operation
	_, ok := err.(gerrors.Gerror)
	if ok {
		return
	}
	errorCode := db.getErrorCode(err)
	err = gerrors.NewFromError(errorCode, err)
	return
}

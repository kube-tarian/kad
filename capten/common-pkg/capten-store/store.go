package captenstore

import (
	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/gerrors"
	postgresdb "github.com/kube-tarian/kad/capten/common-pkg/postgres"
)

type Store struct {
	dbClient *postgresdb.DBClient
	log      logging.Logger
}

func NewStore(log logging.Logger) (*Store, error) {
	dbClient, err := postgresdb.NewDBClient(log)
	if err != nil {
		return nil, err
	}
	return &Store{log: log, dbClient: dbClient}, nil
}

func prepareError(err error, name string, operationLog string) (returnErr error) {
	if gErr, ok := err.(gerrors.Gerror); ok {
		switch gerrors.GetErrorType(gErr) {
		case postgresdb.ObjectNotExist:
			returnErr = gerrors.Newf(gerrors.NotFound, "%s failed for %s not exist", operationLog, name)
		case postgresdb.DuplicateRecord:
			returnErr = gerrors.Newf(gerrors.RecordAlreadyExists, "%s failed, '%s' already exists", operationLog, name)
		case postgresdb.PostgresDBError:
			returnErr = gerrors.Newf(gerrors.InternalError, "%s failed for '%s', Reason: %v", operationLog, name, err.Error())
		}
		return
	}
	return err
}

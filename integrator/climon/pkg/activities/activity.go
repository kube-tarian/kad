package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/kube-tarian/kad/integrator/climon/pkg/db/cassandra"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins"
	workerframework "github.com/kube-tarian/kad/integrator/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/integrator/model"
)

type Activities struct {
}

func (a *Activities) DeploymentActivity(ctx context.Context, req model.RequestPayload) (model.ResponsePayload, error) {
	logger := logging.NewLogger()
	logger.Infof("Activity, name: %+v", req.ToString())
	// e := activity.GetInfo(ctx)
	// logger.Infof("activity info: %+v", e)

	plugin, err := plugins.GetPlugin(req.PluginName, logger)
	if err != nil {
		logger.Errorf("Get plugin  failed: %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}
	deployerPlugin, ok := plugin.(workerframework.DeploymentWorker)
	if !ok {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, fmt.Errorf("plugin not supports deployment activities")
	}

	msg, err := deployerPlugin.DeployActivities(req)
	if err != nil {
		logger.Errorf("Deploy activities failed %s: %v", req.Action, err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	if req.Action == "install" || req.Action == "update" {
		if err := InsertToDb(logger, req.Data); err != nil {
			logger.Errorf("insert db failed, %v", err)
			return model.ResponsePayload{
				Status:  "Failed",
				Message: json.RawMessage(fmt.Sprintf("database update failed %v", err)),
			}, err
		}
	} else if req.Action == "delete" {
		if err := DeleteDbEntry(logger, req.Data); err != nil {
			logger.Errorf("delete plugin failed, %v", err)
			return model.ResponsePayload{
				Status:  "Failed",
				Message: json.RawMessage(fmt.Sprintf("database update failed %v", err)),
			}, err
		}
	}

	return model.ResponsePayload{
		Status:  "Success",
		Message: msg,
	}, nil
}

func InsertToDb(logger logging.Logger, reqData json.RawMessage) error {
	data := &model.Request{}
	if err := json.Unmarshal(reqData, data); err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	dbConf, err := cassandra.GetDbConfig()
	if err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	db, err := cassandra.NewCassandraStore(logger, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	if err := db.InsertToolsDb(data); err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	return nil
}

func DeleteDbEntry(logger logging.Logger, reqData json.RawMessage) error {
	data := &model.Request{}
	if err := json.Unmarshal(reqData, data); err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	dbConf, err := cassandra.GetDbConfig()
	if err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	db, err := cassandra.NewCassandraStore(logger, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	if err := db.DeleteToolsDbEntry(data); err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	return nil
}

package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kube-tarian/kad/integrator/climon/pkg/db/cassandra"
	"github.com/kube-tarian/kad/integrator/common-pkg/logging"
	"github.com/kube-tarian/kad/integrator/common-pkg/plugins"
	workerframework "github.com/kube-tarian/kad/integrator/common-pkg/worker-framework"
	"github.com/kube-tarian/kad/integrator/model"
	"github.com/pkg/errors"
)

type Activities struct {
}

var logger = logging.NewLogger()

func (a *Activities) ClimonInstallActivity(ctx context.Context, req *model.ClimonPostRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)
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

	climonPlugin, ok := plugin.(workerframework.ClimonWorker)
	if !ok {
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage("{\"error\": \"not impmented climon worker plugin\"}"),
		}, fmt.Errorf("plugin not supports deployment activities")
	}

	emptyVersion := ""
	if req.Version == nil {
		req.Version = &emptyVersion
	}
	msg, err := climonPlugin.Create(&model.CreteRequestPayload{
		RepoName:    req.RepoName,
		RepoURL:     req.RepoUrl,
		ChartName:   req.ChartName,
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Timeout:     req.Timeout,
		Version:     *req.Version,
	})
	if err != nil {
		logger.Errorf("Deploy activities failed %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	if err := InsertToDb(logger, req); err != nil {
		logger.Errorf("insert db failed, %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("database update failed %v", err)),
		}, err
	}
	return model.ResponsePayload{
		Status:  "Success",
		Message: msg,
	}, nil
}

func (a *Activities) ClimonDeleteActivity(ctx context.Context, req *model.ClimonDeleteRequest) (model.ResponsePayload, error) {
	logger.Infof("Activity, name: %+v", req)
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
	msg, err := deployerPlugin.Delete(&model.DeleteRequestPayload{
		Namespace:   req.Namespace,
		ReleaseName: req.ReleaseName,
		Timeout:     req.Timeout,
		ClusterName: *req.ClusterName,
	})
	if err != nil {
		logger.Errorf("Deploy activities failed %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("{\"error\": \"%v\"}", strings.ReplaceAll(err.Error(), "\"", "\\\""))),
		}, err
	}

	if err := DeleteDbEntry(logger, req); err != nil {
		logger.Errorf("delete plugin failed, %v", err)
		return model.ResponsePayload{
			Status:  "Failed",
			Message: json.RawMessage(fmt.Sprintf("database update failed %v", err)),
		}, err
	}

	return model.ResponsePayload{
		Status:  "Success",
		Message: msg,
	}, nil
}

func InsertToDb(logger logging.Logger, reqData *model.ClimonPostRequest) error {
	dbConf, err := cassandra.GetDbConfig()
	if err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	db, err := cassandra.NewCassandraStore(logger, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}

	if err := db.InsertToolsDb(reqData); err != nil {
		return errors.Wrap(err, "failed to store data in database")
	}
	return nil
}

func DeleteDbEntry(logger logging.Logger, reqData *model.ClimonDeleteRequest) error {
	dbConf, err := cassandra.GetDbConfig()
	if err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	db, err := cassandra.NewCassandraStore(logger, dbConf.DbAddresses, dbConf.DbAdminUsername, dbConf.DbAdminPassword)
	if err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}

	if err := db.DeleteToolsDbEntry(reqData); err != nil {
		return errors.Wrap(err, "failed to delete data in database")
	}
	return nil
}

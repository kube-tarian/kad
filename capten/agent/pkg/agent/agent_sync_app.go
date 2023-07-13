package agent

import (
	"context"

	"github.com/kube-tarian/kad/capten/agent/pkg/agentpb"
	"github.com/kube-tarian/kad/capten/agent/pkg/types"
	"gopkg.in/yaml.v2"
)

func (a Agent) syncApp(ctx context.Context, request *agentpb.SyncAppRequest) error {

	var appConfig types.AppConfig
	if err := yaml.Unmarshal(request.Payload, &appConfig); err != nil {
		a.log.Errorf("could not unmarshal appConfig yaml: %v", err)
		return err
	}

	// // Note: We are creating keyspace and table when needed
	// if err := a.store.Session().Query(fmt.Sprintf(createKeyspaceCQL, "1")).Exec(); err != nil {
	// 	a.log.Errorf("could not create keyspace, err: %v", err)
	// 	return err
	// }

	// if err := a.store.Session().Query(createTableQuery).Exec(); err != nil {
	// 	a.log.Errorf("could not create table, err: %v", err)
	// 	return err
	// }

	if err := a.Store.InsertAppConfig(appConfig); err != nil {
		a.log.Errorf("could not insert, err: %v", err)
		return err
	}

	return nil
}

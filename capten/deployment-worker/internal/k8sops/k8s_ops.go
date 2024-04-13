package k8sops

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kube-tarian/kad/capten/common-pkg/k8s"
)

func CreateUpdateConfigmap(ctx context.Context, log logging.Logger, namespace, cmName string, data map[string]string, k8sClient *k8s.K8SClient) error {
	err := k8sClient.CreateNamespace(ctx, namespace)
	if err != nil {
		log.Errorf("Creation of namespace failed: %v", err)
		return fmt.Errorf("creation of namespace faield")
	}
	cm, err := k8sClient.GetConfigmap(ctx, namespace, cmName)
	if err != nil {
		log.Infof("plugin configmap %s not found", cmName)
		err = k8sClient.CreateConfigmap(ctx, namespace, cmName, data, map[string]string{})
		if err != nil {
			return fmt.Errorf("failed to create configmap %v", cmName)
		}
	}
	// configmap found but data is empty/nil
	if cm == nil {
		cm = map[string]string{}
	}
	for k, v := range data {
		cm[k] = v
	}
	err = k8sClient.UpdateConfigmap(ctx, namespace, cmName, cm)
	if err != nil {
		return fmt.Errorf("plugin configmap %s not found", cmName)
	}
	return nil
}

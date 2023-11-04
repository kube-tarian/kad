package k8s

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/intelops/go-common/logging"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8SClientController interface {
	Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error
}

type ArgoCDController interface {
	Get(ctx context.Context, name, namespace string) (*v1alpha1.Application, error)
}

type ArgoCDControllerResource struct {
	Namespace string
	Client    K8SClientController
}

func addArgoCDScheme() (*runtime.Scheme, error) {
	newScheme := runtime.NewScheme()

	err := v1alpha1.AddToScheme(newScheme)
	if err != nil {

		return nil, fmt.Errorf("adding argocd schema to runtime instance operation failed, error: %w", err)
	}

	return newScheme, nil
}

func getConfigClient(restConfig *rest.Config) (client.Client, error) {
	newScheme, addErr := addArgoCDScheme()

	if addErr != nil {
		return nil, addErr
	}

	clientObj, err := client.New(restConfig, client.Options{Scheme: newScheme})
	if err != nil {
		return nil, fmt.Errorf("failed to get k8s client object, error: %w", err)
	}

	return clientObj, nil
}

func NewArgoCDClient(logger logging.Logger, namespace string) (ArgoCDController, error) {
	restConfig, err := GetK8SConfig(logger)
	if err != nil {
		return nil, err
	}

	client, err := getConfigClient(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create new argocd custom resource type, error: %w", err)
	}

	return &ArgoCDControllerResource{Namespace: namespace, Client: client}, nil
}

func (cr *ArgoCDControllerResource) Get(ctx context.Context, instanceName, namespace string) (*v1alpha1.Application, error) {

	argoCr := &v1alpha1.Application{}
	// Get call provids an empty instances if resource not found, need to handle explicitly
	gErr := cr.Client.Get(ctx, types.NamespacedName{Namespace: cr.Namespace, Name: instanceName}, argoCr)
	if gErr != nil || argoCr == nil {
		return nil, fmt.Errorf("failed to get the custom resource, error: %w", gErr)
	}

	return argoCr, nil
}

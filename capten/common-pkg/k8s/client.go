package k8s

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Configuration struct {
	KubeconfigPath string `envconfig:"KUBECONFIG_PATH" required:"false"`
}

type K8SClient struct {
	log                    logging.Logger
	Clientset              kubernetes.Interface
	DynamicClientInterface dynamic.Interface
	DynamicClient          *DynamicClientSet
}

func NewK8SClient(log logging.Logger) (*K8SClient, error) {
	config, err := GetK8SConfig(log)
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("initialize kubernetes client failed: %v", err)
	}

	dcClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize dynamic client failed: %v", err)
	}

	return &K8SClient{
		log:                    log,
		Clientset:              clientset,
		DynamicClientInterface: dcClient,
		DynamicClient:          NewDynamicClientSet(dcClient),
	}, nil
}

func GetK8SConfig(log logging.Logger) (*rest.Config, error) {
	conf, err := FetchConfiguration()
	if err != nil {
		return nil, fmt.Errorf("fetch configuration failed: %v", err)
	}

	var k8sConfig *rest.Config
	if conf.KubeconfigPath == "" {
		// creates the in-cluster config
		k8sConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("fetch in-cluster configuration failed: %v", err)
		}
	} else {
		// use the current context in kubeconfig
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", conf.KubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("in-cluster configuration from absolute path %s failed: %v", conf.KubeconfigPath, err)

		}
	}
	return k8sConfig, nil
}

func (k *K8SClient) ListPods(namespace string) ([]corev1.Pod, error) {
	pods, err := k.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func FetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}

func (k *K8SClient) GetSecretData(namespace, secretName string) (*SecretData, error) {
	secret, err := k.Clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		data[k] = string(v)
	}

	return &SecretData{
		Data: data,
	}, nil
}

func (k *K8SClient) GetServiceData(namespace, serviceName string) (*ServiceData, error) {
	service, err := k.Clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Prepare a list of ports
	ports := []int32{}
	for _, v := range service.Spec.Ports {
		ports = append(ports, v.Port)
	}

	return &ServiceData{
		Name:  service.Name,
		Ports: ports,
	}, nil
}

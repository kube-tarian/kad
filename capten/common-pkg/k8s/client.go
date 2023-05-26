package k8s

import (
	"context"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/kube-tarian/kad/capten/common-pkg/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Configuration struct {
	KubeconfigPath string `envconfig:"KUBECONFIG_PATH" required:"false"`
}

type K8SClient struct {
	log       logging.Logger
	Clientset kubernetes.Interface
}

func NewK8SClient(log logging.Logger) (*K8SClient, error) {
	config, err := GetK8SConfig(log)
	if err != nil {
		return nil, err
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Initialize kubernetes client failed: %v", err)
		return nil, err
	}

	return &K8SClient{
		log:       log,
		Clientset: clientset,
	}, nil
}

func GetK8SConfig(log logging.Logger) (*rest.Config, error) {
	conf, err := FetchConfiguration()
	if err != nil {
		log.Errorf("Fetch configuration failed: %v", err)
		return nil, err
	}
	var k8sConfig *rest.Config
	if conf.KubeconfigPath == "" {
		// creates the in-cluster config
		k8sConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Errorf("Fetch in-cluster configuration failed: %v", err)
			return nil, err
		}
	} else {
		// use the current context in kubeconfig
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", conf.KubeconfigPath)
		if err != nil {
			log.Errorf("Fetch in-cluster configuration from absolute path %s failed: %v", conf.KubeconfigPath, err)
			return nil, err
		}
	}
	return k8sConfig, nil
}

func (k *K8SClient) ListPods(namespace string) ([]corev1.Pod, error) {
	pods, err := k.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		k.log.Errorf("List pods failed, %v", err)
		return nil, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	return pods.Items, nil
}

func FetchConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	err := envconfig.Process("", cfg)
	return cfg, err
}

func (k *K8SClient) FetchSecretDetails(req *SecretDetailsRequest) (*SecretDetailsResponse, error) {
	secret, err := k.Clientset.CoreV1().Secrets(req.Namespace).Get(context.TODO(), req.SecretName, metav1.GetOptions{})
	if err != nil {
		k.log.Errorf("Fetching secret %s failed, %v", req.SecretName, err)
		return nil, err
	}

	// Convert data value from []bytes to string
	data := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		data[k] = string(v)
	}

	return &SecretDetailsResponse{
		Namespace: req.Namespace,
		Data:      data,
	}, nil
}

func (k *K8SClient) FetchServiceDetails(req *ServiceDetailsRequest) (*ServiceDetailsResponse, error) {
	service, err := k.Clientset.CoreV1().Services(req.Namespace).Get(context.TODO(), req.ServiceName, metav1.GetOptions{})
	if err != nil {
		k.log.Errorf("Fetching service %s details failed, %v", req.ServiceName, err)
		return nil, err
	}

	// Prepare a list of ports
	ports := []int32{}
	for _, v := range service.Spec.Ports {
		ports = append(ports, v.Port)
	}

	return &ServiceDetailsResponse{
		Namespace: req.Namespace,
		ServiceDetails: ServiceDetails{
			Name:  service.Name,
			Ports: ports,
		},
	}, nil
}

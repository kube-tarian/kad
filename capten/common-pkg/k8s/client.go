package k8s

import (
	"context"
	"fmt"

	"github.com/intelops/go-common/logging"
	"github.com/kelseyhightower/envconfig"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
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
	Config                 *rest.Config
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
		Config:                 config,
	}, nil
}

func NewK8SClientForCluster(log logging.Logger, kubeconfig, clusterCA, endpoint string) (*K8SClient, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		log.Fatal(err.Error())
	}

	config.Host = endpoint
	config.CAData = []byte(clusterCA)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
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
		Config:                 config,
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

func (k *K8SClient) CreateConfigmap(namespace, cmName string, data map[string]string, annotation map[string]string) error {
	_, err := k.Clientset.CoreV1().ConfigMaps(namespace).Create(
		context.TODO(),
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: cmName, Annotations: annotation},
			Data:       data,
		},
		metav1.CreateOptions{})
	return err
}

func (k *K8SClient) UpdateConfigmap(namespace, cmName string, data map[string]string) error {
	_, err := k.Clientset.CoreV1().ConfigMaps(namespace).Update(
		context.TODO(),
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{Name: cmName},
			Data:       data,
		},
		metav1.UpdateOptions{})
	return err
}

func (k *K8SClient) DeleteConfigmap(namespace, cmName string) error {
	cm, _ := k.Clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if cm != nil {
		return k.Clientset.CoreV1().ConfigMaps(namespace).Delete(context.TODO(), cmName, metav1.DeleteOptions{})
	}
	return nil
}

func (k *K8SClient) GetConfigmap(namespace, cmName string) (map[string]string, error) {
	cm, err := k.Clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), cmName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return cm.Data, nil
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

func (k *K8SClient) CreateOrUpdateSecret(ctx context.Context, namespace, secretName string, secretType v1.SecretType,
	data map[string][]byte, annotation map[string]string) error {
	_, err := k.Clientset.CoreV1().Secrets(namespace).Create(ctx,
		&v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName,
			Annotations: annotation},
			Type: secretType, Data: data},
		metav1.CreateOptions{})
	if k8serror.IsAlreadyExists(err) {
		_, err := k.Clientset.CoreV1().Secrets(namespace).Update(ctx,
			&v1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Annotations: annotation},
				Type: secretType, Data: data},
			metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update k8s secret, %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to create k8s secret, %v", err)
	}
	return nil
}

func (k *K8SClient) DeleteSecret(ctx context.Context, namespace, secretName string) error {
	err := k.Clientset.CoreV1().Secrets(namespace).Delete(ctx, secretName, metav1.DeleteOptions{})
	if k8serror.IsNotFound(err) {
		k.log.Info("k8s secret not found %s in namespace", secretName, namespace)
	} else if err != nil {
		return fmt.Errorf("failed to delete k8s secret, %v", err)
	}
	return nil
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

package k8s

type SecretData struct {
	Namespace string
	Data      map[string]string
}

type ServiceData struct {
	Name  string
	Ports []int32
}

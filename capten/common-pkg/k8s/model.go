package k8s

type SecretDetailsRequest struct {
	Namespace  string
	SecretName string
}

type SecretDetailsResponse struct {
	Namespace string
	Data      map[string]string
}

type ServiceDetailsRequest struct {
	Namespace   string
	ServiceName string
}

type ServiceDetails struct {
	Name  string
	Ports []int32
}

type ServiceDetailsResponse struct {
	Namespace string
	ServiceDetails
}

package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/yaml"
)

type DynamicClientSet struct {
	client dynamic.Interface
}

func NewDynamicClientSet(dynamicClient dynamic.Interface) *DynamicClientSet {
	return &DynamicClientSet{client: dynamicClient}
}

func ConvertYamlToJson(data []byte) ([]byte, error) {
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func ConvertJsonToYaml(data []byte) ([]byte, error) {
	yamlData, err := yaml.JSONToYAML(data)
	if err != nil {
		return nil, err
	}

	return yamlData, nil
}

func (dc *DynamicClientSet) GetNameNamespace(jsonByte []byte) (string, string, error) {
	var keyValue map[string]interface{}
	if err := json.Unmarshal(jsonByte, &keyValue); err != nil {
		return "", "", nil
	}

	metadataObj, convCheck := keyValue["metadata"].(map[string]interface{})
	if !convCheck {
		return "", "", fmt.Errorf("failed to convert the metadata togo struct type")
	}

	namespaceName, convCheck := metadataObj["namespace"].(string)
	if !convCheck {
		return "", "", fmt.Errorf("failed to convert the metadata togo struct type")
	}

	resourceName, convCheck := metadataObj["name"].(string)
	if !convCheck {
		return "", "", fmt.Errorf("failed to convert the metadata togo struct type")
	}

	return namespaceName, resourceName, nil
}

func (dc *DynamicClientSet) getGVK(data []byte) (obj *unstructured.Unstructured, resourceID schema.GroupVersionResource, err error) {
	dec := kubeyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	obj = &unstructured.Unstructured{}

	_, gvk, err := dec.Decode([]byte(string(data)), nil, obj)
	if err != nil {
		return
	}

	resourceID = schema.GroupVersionResource{
		Group:    gvk.Group,
		Version:  gvk.Version,
		Resource: strings.ToLower(gvk.Kind + string('s')),
	}

	return
}

func (dc *DynamicClientSet) CreateResourceFromFile(ctx context.Context, filename string) (string, string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", "", err
	}

	return dc.CreateResource(ctx, data)
}

func (dc *DynamicClientSet) CreateResource(ctx context.Context, data []byte) (string, string, error) {
	jsonData, err := ConvertYamlToJson(data)
	if err != nil {
		return "", "", err
	}

	obj, resourceID, err := dc.getGVK(jsonData)
	if err != nil {
		return "", "", err
	}

	namespaceName, resourceName, err := dc.GetNameNamespace(jsonData)
	if err != nil {
		return "", "", err
	}

	_, err = dc.client.Resource(resourceID).Namespace(namespaceName).Get(ctx, resourceName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			_, err := dc.client.Resource(resourceID).Namespace(namespaceName).Create(ctx, obj, metav1.CreateOptions{})
			if err != nil {
				return "", "", fmt.Errorf("error in creating resource %s/%s, %v", namespaceName, resourceName, err)
			}
			return namespaceName, resourceName, nil
		}
		return "", "", err
	}
	return namespaceName, resourceName, nil
}

func (dc *DynamicClientSet) GetResource(ctx context.Context, filename string) (*unstructured.Unstructured, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	jsonData, err := ConvertYamlToJson(data)
	if err != nil {
		return nil, err
	}

	_, resourceID, err := dc.getGVK(jsonData)
	if err != nil {
		return nil, err
	}

	namespaceName, resourceName, err := dc.GetNameNamespace(jsonData)
	if err != nil {
		return nil, err
	}

	obj, err := dc.client.Resource(resourceID).Namespace(namespaceName).Get(ctx, resourceName, metav1.GetOptions{})

	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (dc *DynamicClientSet) ListNamespaceResource(ctx context.Context, gvk schema.GroupVersionResource, ns string) (*unstructured.UnstructuredList, error) {
	objList, err := dc.client.Resource(gvk).Namespace(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return objList, nil
}

func (dc *DynamicClientSet) ListAllNamespaceResource(ctx context.Context, gvk schema.GroupVersionResource) (*unstructured.UnstructuredList, error) {
	objList, err := dc.client.Resource(gvk).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return objList, nil
}

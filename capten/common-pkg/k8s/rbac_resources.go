package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8SClient) CreateOrUpdateServiceAccount(ctx context.Context, namespace, serviceAccountName string) error {
	_, err := k.Clientset.CoreV1().ServiceAccounts(namespace).Create(ctx,
		&v1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: serviceAccountName}},
		metav1.CreateOptions{})
	if !k8serror.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create k8s secret, %v", err)
	}
	return nil
}

func (k *K8SClient) CreateOrUpdateClusterRoleBinding(ctx context.Context, serviceAccounts map[string]string, clusterRole string) error {

	subjects := []rbacv1.Subject{}
	for serviceAccountName, namespace := range serviceAccounts {
		subject := rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      serviceAccountName,
			Namespace: namespace,
		}
		subjects = append(subjects, subject)
	}

	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vault-role-tokenreview-binding",
		},
		Subjects: subjects,
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     clusterRole,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err := k.Clientset.RbacV1().ClusterRoleBindings().Create(ctx,
		clusterRoleBinding, metav1.CreateOptions{})
	if k8serror.IsAlreadyExists(err) {
		_, err := k.Clientset.RbacV1().ClusterRoleBindings().Update(ctx,
			clusterRoleBinding, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("failed to update k8s secret, %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to create k8s secret, %v", err)
	}
	return nil
}

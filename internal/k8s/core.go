package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8s) CreateNamespace(ctx context.Context, name string) error {
	_, err := k.client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "jetbuild",
			},
		},
	}, metav1.CreateOptions{})

	return err
}

func (k *k8s) ListNamespaces(ctx context.Context) (*corev1.NamespaceList, error) {
	return k.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}

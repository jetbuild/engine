package k8s

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8s) CreateDeployment(ctx context.Context, namespace string) error {
	d := v1.Deployment{}

	_, err := k.client.AppsV1().Deployments(namespace).Create(ctx, &d, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

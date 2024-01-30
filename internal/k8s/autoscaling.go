package k8s

import (
	"context"

	v1 "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8s) CreateHPA(ctx context.Context, namespace string) error {
	h := v1.HorizontalPodAutoscaler{}

	_, err := k.client.AutoscalingV1().HorizontalPodAutoscalers(namespace).Create(ctx, &h, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

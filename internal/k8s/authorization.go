package k8s

import (
	"context"
	"fmt"

	authorizationv1 "k8s.io/api/authorization/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *k8s) HasAdminPrivileges(ctx context.Context) (bool, error) {
	auth, err := k.client.AuthorizationV1().SelfSubjectAccessReviews().Create(ctx, &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{
				Namespace:   "*",
				Verb:        "*",
				Group:       "*",
				Version:     "*",
				Resource:    "*",
				Subresource: "*",
				Name:        "*",
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to create self subject access review: %w", err)
	}

	return auth.Status.Allowed, nil
}

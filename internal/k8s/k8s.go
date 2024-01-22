package k8s

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type K8S interface {
	GetConfig() *api.Config
	GetClusterName() string
	HasAdminPrivileges(ctx context.Context) (bool, error)
	CreateNamespace(ctx context.Context, name string) error
}

type k8s struct {
	client      *kubernetes.Clientset
	config      *api.Config
	clusterName string
}

func New(c api.Config) (K8S, error) {
	r, err := clientcmd.NewDefaultClientConfig(c, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create new default client config: %w", err)
	}

	s, err := kubernetes.NewForConfig(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	return &k8s{
		client: s,
	}, nil
}

func NewFromFormFile(f *multipart.FileHeader) (K8S, error) {
	s, err := f.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open kube config file: %w", err)
	}
	defer s.Close()

	p := make([]byte, f.Size)
	if _, err = s.Read(p); err != nil {
		return nil, fmt.Errorf("failed to read kube config file: %w", err)
	}

	c, err := clientcmd.Load(p)
	if err != nil {
		return nil, fmt.Errorf("failed to load kube config file: %w", err)
	}

	if len(c.Contexts) == 0 {
		return nil, errors.New("kube config should have a context")
	}

	if len(c.CurrentContext) == 0 {
		return nil, errors.New("kube config should have current context")
	}

	cx, ok := c.Contexts[c.CurrentContext]
	if !ok {
		return nil, fmt.Errorf("kube config should have '%s' context", c.CurrentContext)
	}

	if len(cx.Cluster) == 0 {
		return nil, fmt.Errorf("kube config '%s' context should have cluster", c.CurrentContext)
	}

	r, err := clientcmd.NewDefaultClientConfig(*c, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create new default client config: %w", err)
	}

	client, err := kubernetes.NewForConfig(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &k8s{
		client:      client,
		config:      c,
		clusterName: cx.Cluster,
	}, nil
}

func (k *k8s) GetConfig() *api.Config {
	return k.config
}

func (k *k8s) GetClusterName() string {
	return k.clusterName
}

package k8s

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8S interface {
	GetConfig() any
	GetClusterName() string
	HasAdminPrivileges(ctx context.Context) (bool, error)
	CreateNamespace(ctx context.Context, name string) error
	ListNamespaces(ctx context.Context) (*corev1.NamespaceList, error)
	CreateDeployment(ctx context.Context, namespace string) error
	CreateHPA(ctx context.Context, namespace string) error
}

type k8s struct {
	client      *kubernetes.Clientset
	config      any
	clusterName string
}

func New(cfg any) (K8S, error) {
	c, err := json.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load kube config file: %w", err)
	}

	cl, err := clientcmd.Load(c)
	if err != nil {
		return nil, fmt.Errorf("failed to load kube config file: %w", err)
	}

	r, err := clientcmd.NewDefaultClientConfig(*cl, nil).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create new default client config: %w", err)
	}

	// TODO: set r.UserAgent

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

	// TODO: set r.UserAgent

	client, err := kubernetes.NewForConfig(r)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	var cfg any
	if err = yaml.Unmarshal(p, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal kube config file: %w", err)
	}

	return &k8s{
		client:      client,
		config:      cfg,
		clusterName: cx.Cluster,
	}, nil
}

func (k *k8s) GetConfig() any {
	return k.config
}

func (k *k8s) GetClusterName() string {
	return k.clusterName
}

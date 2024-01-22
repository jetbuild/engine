package model

import "k8s.io/client-go/tools/clientcmd/api"

type Cluster struct {
	Name   string      `json:"name,omitempty"`
	Config *api.Config `json:"config,omitempty"`
}

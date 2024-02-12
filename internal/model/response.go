package model

import "github.com/jetbuild/engine/pkg/flow"

type ListClustersResponse struct {
	Items []Cluster `json:"items"`
}

type ListClusterNamespacesResponse struct {
	Items []ClusterNamespace `json:"items"`
}

type ListComponentsResponse struct {
	Items []Component `json:"items"`
}

type ListFlowsResponse struct {
	Items []flow.Flow `json:"items"`
}

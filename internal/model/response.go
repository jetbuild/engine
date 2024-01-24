package model

type ListClustersResponse struct {
	Items []Cluster `json:"items"`
}

type ListClusterNamespacesResponse struct {
	Items []ClusterNamespace `json:"items"`
}

type ListComponentsResponse struct {
	Items []Component `json:"items"`
}

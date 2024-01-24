package model

type ListClustersResponse struct {
	Items []Cluster `json:"items"`
}

type ListNamespacesResponse struct {
	Items []Namespace `json:"items"`
}

type ListComponentsResponse struct {
	Items []Component `json:"items"`
}

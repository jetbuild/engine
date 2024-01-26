package model

type Cluster struct {
	Name   string `json:"name,omitempty"`
	Config any    `json:"config,omitempty"`
}

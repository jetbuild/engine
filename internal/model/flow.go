package model

type Flow struct {
	Name       string          `json:"name,omitempty"`
	Components []FlowComponent `json:"components,omitempty"`
	Runners    []FlowRunner    `json:"runners,omitempty"`
}

type FlowComponent struct {
	Key         string                   `json:"key,omitempty"`
	Version     string                   `json:"version,omitempty"`
	Arguments   map[string]any           `json:"arguments,omitempty"`
	Connections *FlowComponentConnection `json:"connections,omitempty"`
}

type FlowComponentConnection struct {
	Targets []uint `json:"targets,omitempty"`
}

type FlowRunner struct {
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Version   string `json:"version,omitempty"`
}

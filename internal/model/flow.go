package model

type Flow struct {
	Name      string        `json:"name,omitempty"`
	Component FlowComponent `json:"component,omitempty"`
	Stages    []Flow        `json:"stages,omitempty"`
}

type FlowComponent struct {
	Key       string         `json:"key,omitempty"`
	Version   string         `json:"version,omitempty"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

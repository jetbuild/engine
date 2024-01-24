package model

type Flow struct {
	Name      string         `json:"name,omitempty"`
	Component string         `json:"component,omitempty"`
	Arguments map[string]any `json:"arguments,omitempty"`
	Stages    []Flow         `json:"stages,omitempty"`
}

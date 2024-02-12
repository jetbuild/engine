package flow

type Flow struct {
	Name       string      `json:"name,omitempty"`
	Components []Component `json:"components,omitempty"`
	Runners    []Runner    `json:"runners,omitempty"`
}

type Component struct {
	Key         string               `json:"key,omitempty"`
	Version     string               `json:"version,omitempty"`
	Arguments   map[string]any       `json:"arguments,omitempty"`
	Connections *ComponentConnection `json:"connections,omitempty"`
}

type ComponentConnection struct {
	Targets []uint `json:"targets,omitempty"`
}

type Runner struct {
	Cluster   string `json:"cluster,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Version   string `json:"version,omitempty"`
}

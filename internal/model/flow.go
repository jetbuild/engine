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

func (f *Flow) GetComponents() []FlowComponent {
	var components []FlowComponent
	components = append(components, f.Component)

	for _, stage := range f.Stages {
		components = append(components, stage.GetComponents()...)
	}

	return components
}

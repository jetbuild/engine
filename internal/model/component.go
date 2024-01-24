package model

import "fmt"

const (
	ComponentArgumentTypeString ComponentArgumentType = "string"
	ComponentArgumentTypeNumber ComponentArgumentType = "number"
	ComponentArgumentTypeBool   ComponentArgumentType = "bool"
)

type Component struct {
	Version     float64             `json:"version,omitempty" yaml:"version"`
	Key         string              `json:"key,omitempty" yaml:"key"`
	Name        string              `json:"name,omitempty" yaml:"name"`
	Description string              `json:"description,omitempty" yaml:"description"`
	Trigger     *bool               `json:"trigger" yaml:"trigger"`
	Arguments   []ComponentArgument `json:"arguments,omitempty" yaml:"arguments"`
}

type ComponentArgument struct {
	Key         string                `json:"key,omitempty" yaml:"key"`
	Name        string                `json:"name,omitempty" yaml:"name"`
	Description string                `json:"description,omitempty" yaml:"description"`
	Type        ComponentArgumentType `json:"type,omitempty" yaml:"type"`
	Required    *bool                 `json:"required" yaml:"required"`
}

type ComponentArgumentType string

func (c *Component) Validate() error {
	if c.Version == 0 {
		return fmt.Errorf("component 'version' field does not found")
	}

	if len(c.Key) == 0 {
		return fmt.Errorf("component 'key' field does not found")
	}

	if len(c.Name) == 0 {
		return fmt.Errorf("component 'name' field does not found")
	}

	if len(c.Description) == 0 {
		return fmt.Errorf("component 'description' field does not found")
	}

	if c.Trigger == nil {
		return fmt.Errorf("component 'trigger' field does not found")
	}

	if len(c.Arguments) == 0 {
		return fmt.Errorf("component 'arguments' field does not found")
	}

	for i, argument := range c.Arguments {
		if len(argument.Key) == 0 {
			return fmt.Errorf("component argument %d 'key' field does not found", i)
		}

		if len(argument.Name) == 0 {
			return fmt.Errorf("component argument %d 'name' field does not found", i)
		}

		if len(argument.Description) == 0 {
			return fmt.Errorf("component argument %d 'description' field does not found", i)
		}

		if len(argument.Type) == 0 {
			return fmt.Errorf("component argument %d 'type' field does not found", i)
		}

		if argument.Type != ComponentArgumentTypeString && argument.Type != ComponentArgumentTypeNumber && argument.Type != ComponentArgumentTypeBool {
			return fmt.Errorf("component argument %d 'type' field does not valid", i)
		}

		if argument.Required == nil {
			return fmt.Errorf("component argument %d 'required' field does not found", i)
		}
	}

	return nil
}

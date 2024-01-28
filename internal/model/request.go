package model

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AddClusterNamespaceRequest struct {
	Body struct {
		Name string `json:"name" validate:"required"`
	}

	Params struct {
		ClusterName string `params:"name" validate:"required"`
	}
}

func (r *AddClusterNamespaceRequest) Bind(ctx *fiber.Ctx, v *validator.Validate) error {
	if err := ctx.BodyParser(&r.Body); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if err := ctx.ParamsParser(&r.Params); err != nil {
		return fmt.Errorf("failed to parse request params: %w", err)
	}

	if err := v.Struct(r.Body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := v.Struct(r.Params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

type ListClusterNamespacesRequest struct {
	Params struct {
		ClusterName string `params:"name" validate:"required"`
	}
}

func (r *ListClusterNamespacesRequest) Bind(ctx *fiber.Ctx, v *validator.Validate) error {
	if err := ctx.ParamsParser(&r.Params); err != nil {
		return fmt.Errorf("failed to parse request params: %w", err)
	}

	if err := v.Struct(r.Params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

type AddFlowRequest struct {
	Name       string                    `json:"name" validate:"required"`
	Components []AddFlowRequestComponent `json:"components" validate:"min=1,dive"`
}

type AddFlowRequestComponent struct {
	Key         string                            `json:"key" validate:"required"`
	Version     string                            `json:"-"`
	Arguments   map[string]any                    `json:"arguments"`
	Connections AddFlowRequestComponentConnection `json:"connections"`
}

type AddFlowRequestComponentConnection struct {
	Sources []uint `json:"sources"`
	Targets []uint `json:"targets"`
}

func (r *AddFlowRequest) Bind(ctx *fiber.Ctx, v *validator.Validate, components []Component) error {
	if err := ctx.BodyParser(&r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if err := v.Struct(r); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := r.Validate(components); err != nil {
		return err
	}

	return nil
}

func (r *AddFlowRequest) Validate(components []Component) error {
	for i, c := range r.Components {
		var component *Component

		for _, cp := range components {
			if c.Key == cp.Key {
				component = &cp

				break
			}
		}

		if component == nil {
			return fmt.Errorf("components[%d].key '%s' does not found", i, c.Key)
		}

		r.Components[i].Version = component.Version

		if i == 0 && !*component.Trigger {
			return fmt.Errorf("components[%d] is not a trigger", i)
		}

		if i != 0 && *component.Trigger {
			return fmt.Errorf("components[%d] cannot be a trigger", i)
		}

		for k, v := range c.Arguments {
			var found *ComponentArgument
			for _, arg := range component.Arguments {
				if k == arg.Key {
					found = &arg

					break
				}
			}

			if found == nil {
				return fmt.Errorf("components[%d].arguments '%s' is not found", i, k)
			}

			if v == nil || v == "" {
				return fmt.Errorf("components[%d].arguments '%s' is empty", i, k)
			}

			if found.Type == ComponentArgumentTypeString && reflect.ValueOf(v).Kind() != reflect.String {
				return fmt.Errorf("components[%d].arguments '%s' is not a string", i, k)
			}

			if found.Type == ComponentArgumentTypeNumber && reflect.ValueOf(v).Kind() != reflect.Float64 {
				return fmt.Errorf("components[%d].arguments '%s' is not a number", i, k)
			}

			if found.Type == ComponentArgumentTypeBool && reflect.ValueOf(v).Kind() != reflect.Bool {
				return fmt.Errorf("components[%d].arguments '%s' is not a bool", i, k)
			}
		}

		for _, arg := range component.Arguments {
			if !*arg.Required {
				continue
			}

			if _, exist := c.Arguments[arg.Key]; !exist {
				return fmt.Errorf("components[%d].arguments '%s' is required", i, arg.Key)
			}
		}

		// TODO: add connections validation
	}

	return nil
}

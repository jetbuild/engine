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
	Name      string           `json:"name" validate:"required"`
	Component string           `json:"component" validate:"required"`
	Arguments map[string]any   `json:"arguments"`
	Stages    []AddFlowRequest `json:"stages" validate:"dive"`
}

func (r *AddFlowRequest) Bind(ctx *fiber.Ctx, v *validator.Validate, components []Component) error {
	if err := ctx.BodyParser(&r); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if err := v.Struct(r); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := r.Validate(true, components); err != nil {
		return err
	}

	return nil
}

func (r *AddFlowRequest) Validate(trigger bool, components []Component) error {
	var component *Component
	for _, c := range components {
		if c.Key == r.Component {
			component = &c

			break
		}
	}

	if component == nil {
		return fmt.Errorf("component '%s' does not found", r.Component)
	}

	if trigger != *component.Trigger {
		return fmt.Errorf("component '%s' trigger type is invalid", r.Component)
	}

	for k, v := range r.Arguments {
		var found *ComponentArgument
		for _, arg := range component.Arguments {
			if k == arg.Key {
				found = &arg

				break
			}
		}

		if found == nil {
			return fmt.Errorf("argument '%s' is not found", k)
		}

		if v == nil || v == "" {
			return fmt.Errorf("argument '%s' is empty", k)
		}

		if found.Type == ComponentArgumentTypeString && reflect.ValueOf(v).Kind() != reflect.String {
			return fmt.Errorf("argument '%s' is not a string", k)
		}

		if found.Type == ComponentArgumentTypeNumber && reflect.ValueOf(v).Kind() != reflect.Float64 {
			return fmt.Errorf("argument '%s' is not a number", k)
		}

		if found.Type == ComponentArgumentTypeBool && reflect.ValueOf(v).Kind() != reflect.Bool {
			return fmt.Errorf("argument '%s' is not a bool", k)
		}
	}

	for _, arg := range component.Arguments {
		if !*arg.Required {
			continue
		}

		if _, exist := r.Arguments[arg.Key]; !exist {
			return fmt.Errorf("argument '%s' is required", arg.Key)
		}
	}

	for _, s := range r.Stages {
		if err := s.Validate(false, components); err != nil {
			return err
		}
	}

	return nil
}

func (r *AddFlowRequest) ToFlow() Flow {
	flow := Flow{
		Name:      r.Name,
		Component: r.Component,
		Arguments: r.Arguments,
	}

	for _, stage := range r.Stages {
		flow.Stages = append(flow.Stages, stage.ToFlow())
	}

	return flow
}

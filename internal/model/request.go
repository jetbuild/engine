package model

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AddNamespaceRequest struct {
	Body struct {
		Name string `json:"name" validate:"required"`
	}

	Params struct {
		ClusterName string `params:"name" validate:"required"`
	}
}

func (r *AddNamespaceRequest) Bind(ctx *fiber.Ctx, validator *validator.Validate) error {
	if err := ctx.BodyParser(&r.Body); err != nil {
		return fmt.Errorf("failed to parse request body: %w", err)
	}

	if err := ctx.ParamsParser(&r.Params); err != nil {
		return fmt.Errorf("failed to parse request params: %w", err)
	}

	if err := validator.Struct(&r.Body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err := validator.Struct(&r.Params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

type ListNamespacesRequest struct {
	Params struct {
		ClusterName string `params:"name" validate:"required"`
	}
}

func (r *ListNamespacesRequest) Bind(ctx *fiber.Ctx, validator *validator.Validate) error {
	if err := ctx.ParamsParser(&r.Params); err != nil {
		return fmt.Errorf("failed to parse request params: %w", err)
	}

	if err := validator.Struct(&r.Params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return nil
}

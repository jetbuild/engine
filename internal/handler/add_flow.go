package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"github.com/jetbuild/engine/pkg/flow"
)

func (h *Handler) addFlow(ctx *fiber.Ctx) error {
	var req model.AddFlowRequest
	if err := req.Bind(ctx, h.Validator, h.Components); err != nil {
		return err
	}

	f := flow.Flow{
		Name: req.Name,
	}

	for _, c := range req.Components {
		f.Components = append(f.Components, flow.Component{
			Key:       c.Key,
			Version:   c.Version,
			Arguments: c.Arguments,
			Connections: &flow.ComponentConnection{
				Targets: c.Connections.Targets,
			},
		})
	}

	err := h.FlowRepository.Add(ctx.Context(), req.Name, f)
	if err != nil && errors.Is(err, vault.ErrItemAlreadyExist) {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("flow '%s' already exist", req.Name))
	}
	if err != nil {
		return fmt.Errorf("failed to save flow to vault: %w", err)
	}

	ctx.Status(fiber.StatusCreated)

	return nil
}

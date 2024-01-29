package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) addFlow(ctx *fiber.Ctx) error {
	var req model.AddFlowRequest
	if err := req.Bind(ctx, h.Validator, h.Components); err != nil {
		return err
	}

	f := model.Flow{
		Name: req.Name,
	}

	for _, c := range req.Components {
		f.Components = append(f.Components, model.FlowComponent{
			Key:       c.Key,
			Version:   c.Version,
			Arguments: c.Arguments,
			Connections: &model.FlowComponentConnection{
				Targets: c.Connections.Targets,
			},
		})
	}

	err := h.FlowRepository.Add(ctx.UserContext(), req.Name, f)
	if err != nil && errors.Is(err, vault.ErrItemAlreadyExist) {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("flow '%s' already exist", req.Name))
	}
	if err != nil {
		return fmt.Errorf("failed to save flow to vault: %w", err)
	}

	ctx.Status(fiber.StatusCreated)

	return nil
}

package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) listFlows(ctx *fiber.Ctx) error {
	res := model.ListFlowsResponse{
		Items: make([]model.Flow, 0),
	}

	flows, err := h.FlowRepository.List(ctx.Context())
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(res)
	}
	if err != nil {
		return fmt.Errorf("failed to list flows: %w", err)
	}

	for _, flow := range flows {
		res.Items = append(res.Items, flow)
	}

	return ctx.JSON(res)
}

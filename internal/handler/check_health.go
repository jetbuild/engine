package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) checkHealth(ctx *fiber.Ctx) error {
	if err := h.ClusterRepository.Ping(ctx.UserContext()); err != nil {
		return fmt.Errorf("failed to ping vault: %w", err)
	}

	ctx.Status(fiber.StatusNoContent)

	return nil
}

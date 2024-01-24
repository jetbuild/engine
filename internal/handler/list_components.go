package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
)

func (h *Handler) listComponents(ctx *fiber.Ctx) error {
	return ctx.JSON(model.ListComponentsResponse{
		Items: h.Components,
	})
}

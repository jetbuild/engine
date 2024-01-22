package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) listClusters(ctx *fiber.Ctx) error {
	res := model.ListClustersResponse{
		Items: make([]model.Cluster, 0),
	}

	clusters, err := h.ClusterRepository.List(ctx.UserContext())
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return ctx.Status(fiber.StatusNotFound).JSON(res)
	}
	if err != nil {
		return fmt.Errorf("failed to list clusters: %w", err)
	}

	for _, cluster := range clusters {
		res.Items = append(res.Items, model.Cluster{
			Name: cluster.Name,
		})
	}

	return ctx.JSON(res)
}

package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/k8s"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) listClusterNamespaces(ctx *fiber.Ctx) error {
	var req model.ListClusterNamespacesRequest
	if err := req.Bind(ctx, h.Validator); err != nil {
		return err
	}

	res := model.ListClusterNamespacesResponse{
		Items: make([]model.ClusterNamespace, 0),
	}

	cluster, err := h.ClusterRepository.Get(ctx.Context(), req.Params.ClusterName)
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "cluster does not found in vault")
	}
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	c, err := k8s.New(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	namespaces, err := c.ListNamespaces(ctx.Context())
	if err != nil {
		return fmt.Errorf("failed to list cluster namespaces: %w", err)
	}

	for _, namespace := range namespaces.Items {
		res.Items = append(res.Items, model.ClusterNamespace{
			Name: namespace.Name,
		})
	}

	return ctx.JSON(res)
}

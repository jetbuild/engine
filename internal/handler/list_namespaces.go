package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/k8s"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) listNamespaces(ctx *fiber.Ctx) error {
	var req model.ListNamespacesRequest
	if err := req.Bind(ctx, h.Validator); err != nil {
		return err
	}

	res := model.ListNamespacesResponse{
		Items: make([]model.Namespace, 0),
	}

	cluster, err := h.ClusterRepository.Get(ctx.UserContext(), req.Params.ClusterName)
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "cluster does not found in vault")
	}
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	c, err := k8s.New(*cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	namespaces, err := c.ListNamespaces(ctx.UserContext())
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, namespace := range namespaces.Items {
		res.Items = append(res.Items, model.Namespace{
			Name: namespace.Name,
		})
	}

	return ctx.JSON(res)
}

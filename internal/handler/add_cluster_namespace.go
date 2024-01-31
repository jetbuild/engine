package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/k8s"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

func (h *Handler) addClusterNamespace(ctx *fiber.Ctx) error {
	var req model.AddClusterNamespaceRequest
	if err := req.Bind(ctx, h.Validator); err != nil {
		return err
	}

	cluster, err := h.ClusterRepository.Get(ctx.Context(), req.Params.ClusterName)
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "cluster does not found in vault")
	}
	if err != nil {
		return fmt.Errorf("failed to get cluster from vault: %w", err)
	}

	c, err := k8s.New(cluster.Config)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	err = c.CreateNamespace(ctx.Context(), req.Body.Name)
	if apierrors.IsAlreadyExists(err) {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("cluster namespace '%s' already exist", req.Body.Name))
	}
	if err != nil {
		return fmt.Errorf("failed to create cluster namespace: %w", err)
	}

	ctx.Status(fiber.StatusCreated)

	return nil
}

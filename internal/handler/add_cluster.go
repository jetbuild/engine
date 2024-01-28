package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/k8s"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func (h *Handler) addCluster(ctx *fiber.Ctx) error {
	f, err := ctx.FormFile("kubeConfig")
	if err != nil {
		return fmt.Errorf("failed to get kube config file: %w", err)
	}

	cluster, err := k8s.NewFromFormFile(f)
	if err != nil {
		return fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	admin, err := cluster.HasAdminPrivileges(ctx.UserContext())
	if err != nil {
		return fmt.Errorf("failed to check admin priveleges: %w", err)
	}

	if !admin {
		return fiber.NewError(fiber.StatusUnauthorized, "kube config does not have admin privileges")
	}

	name := cluster.GetClusterName()

	err = h.ClusterRepository.Add(ctx.UserContext(), name, model.Cluster{
		Name:   name,
		Config: cluster.GetConfig(),
	})
	if err != nil && errors.Is(err, vault.ErrItemAlreadyExist) {
		return fiber.NewError(fiber.StatusConflict, fmt.Sprintf("cluster '%s' already exist", name))
	}
	if err != nil {
		return fmt.Errorf("failed to save cluster to vault: %w", err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

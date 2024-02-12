package handler

import (
	"errors"
	"fmt"
	"slices"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/k8s"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"github.com/jetbuild/engine/pkg/flow"
	v1 "k8s.io/api/core/v1"
)

func (h *Handler) addFlowRunner(ctx *fiber.Ctx) error {
	var req model.AddFlowRunnerRequest
	if err := req.Bind(ctx, h.Validator); err != nil {
		return err
	}

	f, err := h.FlowRepository.Get(ctx.Context(), req.Params.FlowName)
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "flow does not found in vault")
	}
	if err != nil {
		return fmt.Errorf("failed to get flow from vault: %w", err)
	}

	if slices.ContainsFunc(f.Runners, func(r flow.Runner) bool {
		return r.Cluster == req.Body.Cluster
	}) {
		return fiber.NewError(fiber.StatusConflict, "runner already exist for flow")
	}

	cluster, err := h.ClusterRepository.Get(ctx.Context(), req.Body.Cluster)
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

	if !slices.ContainsFunc(namespaces.Items, func(n v1.Namespace) bool {
		if n.Name == req.Body.Namespace {
			return true
		}

		return false
	}) {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("namespace does not exist in cluster '%s'", req.Body.Cluster))
	}

	// TODO: create necessary k8s resources
	/*
		if err = c.CreateDeployment(ctx.Context(), req.Body.Namespace); err != nil {
			return fmt.Errorf("failed to create deployment: %w", err)
		}

		if err = c.CreateHPA(ctx.Context(), req.Body.Namespace); err != nil {
			return fmt.Errorf("failed to create hpa: %w", err)
		}
	*/

	f.Runners = append(f.Runners, flow.Runner{
		Cluster:   req.Body.Cluster,
		Namespace: req.Body.Namespace,
		Version:   h.LatestRunnerVersion,
	})

	err = h.FlowRepository.Update(ctx.Context(), req.Params.FlowName, *f)
	if err != nil && errors.Is(err, vault.ErrKeyNotFound) {
		return fiber.NewError(fiber.StatusNotFound, "flow does not found in vault for update")
	}
	if err != nil {
		return fmt.Errorf("failed to update flow from vault: %w", err)
	}

	ctx.Status(fiber.StatusCreated)

	return nil
}

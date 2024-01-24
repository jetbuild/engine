package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jetbuild/engine/internal/component"
	"github.com/jetbuild/engine/internal/config"
	"github.com/jetbuild/engine/internal/handler"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	var c config.Config
	if err := c.Load(); err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	timeout, err := time.ParseDuration(c.ServerInitTimeout)
	if err != nil {
		logger.Error("failed to parse server init timeout duration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	v, err := vault.New(ctx, c.VaultAddr, c.VaultToken, c.VaultEngine, c.VaultEngineDescription)
	if err != nil {
		logger.Error("failed to connect vault", "error", err)
		os.Exit(1)
	}

	components, err := component.Load(ctx, c.GithubOrganization)
	if err != nil {
		logger.Error("failed to load components", "error", err)
		os.Exit(1)
	}

	h := handler.Handler{
		Validator:         validator.New(validator.WithRequiredStructEnabled()),
		ClusterRepository: vault.NewRepository[model.Cluster](v, "clusters"),
		FlowRepository:    vault.NewRepository[model.Flow](v, "flows"),
		Config:            &c,
		Components:        components,
	}

	if err = h.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

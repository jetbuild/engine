package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jetbuild/engine/internal/config"
	"github.com/jetbuild/engine/internal/github"
	"github.com/jetbuild/engine/internal/handler"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"github.com/jetbuild/engine/pkg/flow"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	var c config.Config
	if err := c.Load(); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	timeout, err := time.ParseDuration(c.ServerInitTimeout)
	if err != nil {
		slog.Error("failed to parse server init timeout duration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	v, err := vault.New(ctx, c.VaultAddr, c.VaultToken, c.VaultEngine, c.VaultEngineDescription)
	if err != nil {
		slog.Error("failed to connect vault", "error", err)
		os.Exit(1)
	}

	h := handler.Handler{
		Validator:         validator.New(validator.WithRequiredStructEnabled()),
		ClusterRepository: vault.NewRepository[model.Cluster](v, "clusters"),
		FlowRepository:    vault.NewRepository[flow.Flow](v, "flows"),
		Config:            &c,
		GitHub:            github.New(c.GithubOrganization),
	}

	t, err := h.GitHub.GetRepositoryLatestTag(ctx, "runner")
	if err != nil {
		slog.Error("failed to get latest runner repository tag", "error", err)
		os.Exit(1)
	}

	h.LatestRunnerVersion = strings.TrimLeft(t, "v")

	if err = h.Start(); err != nil {
		slog.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

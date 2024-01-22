package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
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

	v, err := vault.New(c.VaultAddr, c.VaultToken, c.VaultEngine, "JetBuild Secret Storage")
	if err != nil {
		logger.Error("failed to connect vault", "error", err)
		os.Exit(1)
	}

	h := handler.Handler{
		Validator:         validator.New(validator.WithRequiredStructEnabled()),
		ClusterRepository: vault.NewRepository[model.Cluster](v, "clusters"),
		Config:            &c,
	}

	if err = h.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

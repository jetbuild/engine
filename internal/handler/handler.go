package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jetbuild/engine/internal/config"
	"github.com/jetbuild/engine/internal/github"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	Validator         *validator.Validate
	ClusterRepository vault.Vault[model.Cluster]
	FlowRepository    vault.Vault[model.Flow]
	Config            *config.Config
	Components        []model.Component
	Logger            *slog.Logger
	GitHub            github.GitHub
}

func (h *Handler) Start() error {
	f := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	f.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	})).Group(h.Config.ServerRoutePrefix).
		Get("/livez", h.checkHealth).
		Get("/readyz", h.checkHealth).
		Get("/clusters", h.listClusters).
		Post("/clusters", h.addCluster).
		Get("/clusters/:name/namespaces", h.listClusterNamespaces).
		Post("/clusters/:name/namespaces", h.addClusterNamespace).
		Get("/components", h.listComponents).
		Get("/flows", h.listFlows).
		Post("/flows", h.addFlow)

	h.Logger.Info("server started",
		slog.String("addr", h.Config.ServerAddr),
		slog.Uint64("handlers", uint64(f.HandlersCount())),
		slog.Int("pid", os.Getpid()),
	)

	ctx := f.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("loadComponents", true)

	if err := h.listComponents(ctx); err != nil {
		return fmt.Errorf("failed to load components: %w", err)
	}

	return f.Listen(h.Config.ServerAddr)
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "an error occurred"

	if err != nil {
		message = err.Error()
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}

	return ctx.Status(code).JSON(struct {
		// TODO: title
		Message string `json:"message"`
	}{
		Message: message,
	})
}

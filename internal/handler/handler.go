package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jetbuild/engine/internal/config"
	"github.com/jetbuild/engine/internal/github"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"github.com/jetbuild/engine/pkg/flow"
	"github.com/valyala/fasthttp"
)

type Handler struct {
	Validator           *validator.Validate
	ClusterRepository   vault.Vault[model.Cluster]
	FlowRepository      vault.Vault[flow.Flow]
	Config              *config.Config
	Components          []model.Component
	GitHub              github.GitHub
	LatestRunnerVersion string
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
		Post("/flows", h.addFlow).
		Post("/flows/:name/runners", h.addFlowRunner)

	f.Hooks().OnListen(func(d fiber.ListenData) error {
		if fiber.IsChild() {
			return nil
		}

		scheme := "http"
		if d.TLS {
			scheme = "https"
		}

		slog.Info("server started",
			slog.String("addr", fmt.Sprintf("%s://%s:%s", scheme, d.Host, d.Port)),
			slog.Uint64("handlers", uint64(f.HandlersCount())),
			slog.Int("pid", os.Getpid()),
		)

		return nil
	})

	f.Hooks().OnShutdown(func() error {
		slog.Info("server stopped")

		return nil
	})

	ctx := f.AcquireCtx(&fasthttp.RequestCtx{})
	ctx.Locals("loadComponents", true)

	if err := h.listComponents(ctx); err != nil {
		return fmt.Errorf("failed to load components: %w", err)
	}

	go func() {
		if err := f.Listen(h.Config.ServerAddr); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	_ = <-c
	slog.Info("server gracefully shutting down")

	return f.Shutdown()
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

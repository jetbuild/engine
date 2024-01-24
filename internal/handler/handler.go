package handler

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jetbuild/engine/internal/config"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
)

type Handler struct {
	Validator         *validator.Validate
	ClusterRepository vault.Vault[model.Cluster]
	FlowRepository    vault.Vault[model.Flow]
	Config            *config.Config
	Components        []model.Component
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
		Get("/clusters/:name/namespaces", h.listNamespaces).
		Post("/clusters/:name/namespaces", h.addNamespace).
		Get("/components", h.listComponents).
		Post("/flows", h.addFlow)

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

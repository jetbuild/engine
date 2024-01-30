package handler

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jetbuild/engine/internal/model"
	"github.com/jetbuild/engine/internal/vault"
	"gopkg.in/yaml.v3"
)

func (h *Handler) listComponents(ctx *fiber.Ctx) error {
	if ctx.Locals("loadComponents") == nil {
		return ctx.JSON(model.ListComponentsResponse{
			Items: h.Components,
		})
	}

	flows, err := h.FlowRepository.List(ctx.UserContext())
	if err != nil && !errors.Is(err, vault.ErrKeyNotFound) {
		return fmt.Errorf("failed to list flows from vault: %w", err)
	}

	components := make(map[string]struct{})
	for _, f := range flows {
		for _, c := range f.Components {
			components[fmt.Sprintf("%s-component:v%s", c.Key, c.Version)] = struct{}{}
		}
	}

	org := h.GitHub.GetOrganizationName()

	repos, err := h.GitHub.ListRepositories(ctx.UserContext())
	if err != nil {
		return fmt.Errorf("failed to list github component repositories by '%s' org: %w", org, err)
	}

	for _, repo := range repos {
		if !strings.HasSuffix(repo.GetName(), "-component") {
			continue
		}

		if !slices.Contains(repo.Topics, fmt.Sprintf("%s-component", org)) {
			continue
		}

		components[fmt.Sprintf("%s:%s", repo.GetName(), "main")] = struct{}{}
	}

	if len(components) == 0 {
		return fmt.Errorf("could not find a component repository on '%s' github org", org)
	}

	var list []model.Component

	for repo := range components {
		s := strings.Split(repo, ":")

		c, cErr := h.GitHub.GetRepositoryContent(ctx.UserContext(), s[0], s[1], "spec.yml")
		if cErr != nil {
			return fmt.Errorf("failed to get component spec file content from github '%s' org '%s' repository: %w", org, s[0], cErr)
		}

		var component model.Component

		if err = yaml.NewDecoder(c).Decode(&component); err != nil {
			return fmt.Errorf("failed to decode component spec file content from github '%s' org '%s' repository: %w", org, s[0], err)
		}

		if err = component.Validate(); err != nil {
			return fmt.Errorf("failed to validate component spec file content from github '%s' org '%s' repository: %w", org, s[0], err)
		}

		list = append(list, component)
	}

	versionMap := make(map[string]bool)

	var f []model.Component

	for _, c := range list {
		if _, ok := versionMap[c.Key+c.Version]; !ok {
			versionMap[c.Key+c.Version] = true
			f = append(f, c)
		}
	}

	h.Components = f

	return nil
}

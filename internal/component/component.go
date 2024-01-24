package component

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/google/go-github/v58/github"
	"github.com/jetbuild/engine/internal/model"

	"gopkg.in/yaml.v3"
)

func Load(ctx context.Context, githubOrganization string) ([]model.Component, error) {
	gh := github.NewClient(nil)

	o := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 1000},
	}

	var repos []*github.Repository
	for {
		l, res, err := gh.Repositories.ListByOrg(ctx, githubOrganization, o)
		if err != nil {
			return nil, fmt.Errorf("failed to list component repositories: %w", err)
		}

		for _, repo := range l {
			if !strings.HasSuffix(repo.GetName(), "-component") {
				continue
			}

			if !slices.Contains(repo.Topics, fmt.Sprintf("%s-component", githubOrganization)) {
				continue
			}

			repos = append(repos, repo)
		}

		if res.NextPage == 0 {
			break
		}

		o.Page = res.NextPage
	}

	if len(repos) == 0 {
		return nil, fmt.Errorf("could not find a component repository on github")
	}

	var components []model.Component

	for _, repo := range repos {
		f, _, _, err := gh.Repositories.GetContents(ctx, githubOrganization, repo.GetName(), "spec.yml", &github.RepositoryContentGetOptions{
			Ref: "main",
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get '%s' component repository spec file: %w", repo.GetName(), err)
		}

		c, err := f.GetContent()
		if err != nil {
			return nil, fmt.Errorf("failed to get '%s' component repository spec file content: %w", repo.GetName(), err)
		}

		var cp model.Component

		if err = yaml.NewDecoder(strings.NewReader(c)).Decode(&cp); err != nil {
			return nil, fmt.Errorf("failed to decode '%s' component spec file: %w", repo.GetName(), err)
		}

		if err = cp.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate '%s' component spec file: %w", repo.GetName(), err)
		}

		components = append(components, cp)
	}

	return components, nil
}

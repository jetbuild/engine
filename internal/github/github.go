package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v58/github"
)

type GitHub interface {
	GetOrganizationName() string
	ListRepositories(ctx context.Context) ([]github.Repository, error)
	GetRepositoryContent(ctx context.Context, name, ref, path string) (*strings.Reader, error)
	GetRepositoryLatestTag(ctx context.Context, name string) (string, error)
}

type gitHub struct {
	client *github.Client
	org    string
}

func New(org string) GitHub {
	return &gitHub{
		client: github.NewClient(nil),
		org:    org,
	}
}

func (g *gitHub) GetOrganizationName() string {
	return g.org
}

func (g *gitHub) ListRepositories(ctx context.Context) ([]github.Repository, error) {
	opt := github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 1000},
	}

	var repos []github.Repository
	for {
		list, res, err := g.client.Repositories.ListByOrg(ctx, g.org, &opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		for _, repo := range list {
			repos = append(repos, *repo)
		}

		if res.NextPage == 0 {
			break
		}

		opt.Page = res.NextPage
	}

	return repos, nil
}

func (g *gitHub) GetRepositoryContent(ctx context.Context, name, ref, path string) (*strings.Reader, error) {
	f, _, _, err := g.client.Repositories.GetContents(ctx, g.org, name, path, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get repository file: %w", err)
	}

	c, err := f.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get repository file content: %w", err)
	}

	return strings.NewReader(c), nil
}

func (g *gitHub) GetRepositoryLatestTag(ctx context.Context, name string) (string, error) {
	r, _, err := g.client.Repositories.GetLatestRelease(ctx, g.org, name)
	if err != nil {
		return "", fmt.Errorf("failed to get repository latest release: %w", err)
	}

	return *r.TagName, nil
}

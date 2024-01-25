package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v58/github"
)

type GitHub interface {
	GetOrganizationName() string
	ListRepositoriesByOrg(ctx context.Context) ([]github.Repository, error)
	GetContentByRepository(ctx context.Context, name, ref, path string) (*strings.Reader, error)
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

func (g *gitHub) ListRepositoriesByOrg(ctx context.Context) ([]github.Repository, error) {
	opt := github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 1000},
	}

	var repos []github.Repository
	for {
		list, res, err := g.client.Repositories.ListByOrg(ctx, g.org, &opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories by org: %w", err)
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

func (g *gitHub) GetContentByRepository(ctx context.Context, name, ref, path string) (*strings.Reader, error) {
	f, _, _, err := g.client.Repositories.GetContents(ctx, g.org, name, path, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get file by repository: %w", err)
	}

	c, err := f.GetContent()
	if err != nil {
		return nil, fmt.Errorf("failed to get content by repository: %w", err)
	}

	return strings.NewReader(c), nil
}

package context

import (
	"fmt"
	"strings"

	"github.com/cli/cli/git"
	"github.com/cli/cli/internal/config"
	"github.com/cli/cli/internal/ghrepo"
)

// NewBlank initializes a blank Context suitable for testing
func NewBlank() *blankContext {
	return &blankContext{}
}

// A Context implementation that queries the filesystem
type blankContext struct {
	branch   string
	baseRepo ghrepo.Interface
	remotes  Remotes
}

func (c *blankContext) Config() (config.Config, error) {
	cfg, err := config.ParseConfig("config.yml")
	if err != nil {
		panic(fmt.Sprintf("failed to parse config during tests. did you remember to stub? error: %s", err))
	}
	return cfg, nil
}

func (c *blankContext) Branch() (string, error) {
	if c.branch == "" {
		return "", fmt.Errorf("branch was not initialized: %w", git.ErrNotOnAnyBranch)
	}
	return c.branch, nil
}

func (c *blankContext) SetBranch(b string) {
	c.branch = b
}

func (c *blankContext) Remotes() (Remotes, error) {
	if c.remotes == nil {
		return nil, fmt.Errorf("remotes were not initialized")
	}
	return c.remotes, nil
}

func (c *blankContext) SetRemotes(stubs map[string]string) {
	c.remotes = make([]*Remote, 0, len(stubs))
	for remoteName, repo := range stubs {
		ownerWithName := strings.SplitN(repo, "/", 2)
		c.remotes = append(c.remotes, &Remote{
			Remote: &git.Remote{Name: remoteName},
			Repo:   ghrepo.New(ownerWithName[0], ownerWithName[1]),
		})
	}
}

func (c *blankContext) BaseRepo() (ghrepo.Interface, error) {
	if c.baseRepo != nil {
		return c.baseRepo, nil
	}
	remotes, err := c.Remotes()
	if err != nil {
		return nil, err
	}
	if len(remotes) < 1 {
		return nil, fmt.Errorf("remotes are empty")
	}
	return remotes[0], nil
}

func (c *blankContext) SetBaseRepo(nwo string) {
	repo, _ := ghrepo.FromFullName(nwo)
	c.baseRepo = repo
}

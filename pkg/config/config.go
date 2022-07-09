package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/github"
	gogit "github.com/go-git/go-git/v5"
	gogittransport "github.com/go-git/go-git/v5/plumbing/transport"
	gogithttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	gogitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	gh "github.com/google/go-github/v44/github"
	"github.com/sethvargo/go-githubactions"
	"golang.org/x/oauth2"
)

type Config struct {
	Trunk         string
	EventSha      string
	ShouldSetTag  bool
	ShouldPushTag bool
	Labelmap      map[string]github.SemVerBump
	Repoclient    *github.RepoClient
	Repo          *gogit.Repository
	RepoAuth      gogittransport.AuthMethod
}

// TODO since you can specify defaults in action.yml, this here is kind of a
// TODO duplication, and should be removed from here
func ConfigInsideActions() (*Config, error) {
	trunk := os.Getenv("GITHUB_REF_NAME")
	if trunk == "" {
		return nil, fmt.Errorf("could not read trunk, env variable GITHUB_REF_NAME is empty")
	}

	shouldSetTag := true
	if githubactions.GetInput("should-set-tag") == "false" {
		shouldSetTag = false
	}

	shouldPushTag := true
	if githubactions.GetInput("should-push-tag") == "false" {
		shouldPushTag = false
	}

	// since labelmap is optional, use the default values if necessary
	major := githubactions.GetInput("label-major")
	if major == "" {
		major = "merge-major"
	}
	minor := githubactions.GetInput("label-minor")
	if minor == "" {
		minor = "merge-minor"
	}
	patch := githubactions.GetInput("label-patch")
	if patch == "" {
		patch = "merge-patch"
	}
	none := githubactions.GetInput("label-none")
	if none == "" {
		none = "merge-none"
	}
	lblmap := map[string]github.SemVerBump{
		major: github.Major,
		minor: github.Minor,
		patch: github.Patch,
		none:  github.None,
	}

	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	if owner == "" {
		return nil, fmt.Errorf("could not read owner, env variable GITHUB_REPOSITORY_OWNER is empty")
	}
	// for some reason the name of a repository is not a dedicated variable, so we need to split it out
	repoName := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")[1]
	if repoName == "" {
		return nil, fmt.Errorf("could not read repoName, env variable GITHUB_REPOSITORY is empty")
	}

	token := githubactions.GetInput("repo-token")
	if token == "" {
		return nil, fmt.Errorf("input variable 'repo-token' cannot be empty")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := gh.NewClient(tc)

	workspace := os.Getenv("GITHUB_WORKSPACE")
	if workspace == "" {
		return nil, fmt.Errorf("could not read workspace, env variable GITHUB_WORKSPACE is empty")
	}
	repoPath := githubactions.GetInput("repo-storage-path-overwrite")
	if repoPath == "" {
		repoPath = workspace
	}

	eventSha := os.Getenv("GITHUB_SHA")
	if eventSha == "" {
		return nil, fmt.Errorf("could not read eventSha, env variable GITHUB_SHA is empty")
	}

	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	var repoAuth gogittransport.AuthMethod
	repoSSHKey := githubactions.GetInput("repo-ssh-key")
	if repoSSHKey == "" {
		// use the "repo-token" as git authentication
		repoAuth = &gogithttp.BasicAuth{
			Username: "githubactions@email.com", // this can be anything except empty, when using with a token
			Password: token,
		}
	} else {
		repoAuth, err = gogitssh.NewPublicKeys("git", []byte(repoSSHKey), "")
		if err != nil {
			return nil, err
		}
	}

	conf := &Config{
		Trunk:         trunk,
		EventSha:      eventSha,
		ShouldSetTag:  shouldSetTag,
		ShouldPushTag: shouldPushTag,
		Labelmap:      lblmap,
		Repoclient: &github.RepoClient{
			Owner:    owner,
			RepoName: repoName,
			// TODO investigate literal copy of lock values
			Client: *client,
		},
		Repo:     repo,
		RepoAuth: repoAuth,
	}

	return conf, nil
}

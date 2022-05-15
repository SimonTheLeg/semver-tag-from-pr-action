package github

import (
	"context"
	"fmt"

	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
	gh "github.com/google/go-github/v44/github"
)

type GitHubRepoClient struct {
	Owner    string
	RepoName string
	gh.Client
}

// GetPRForCommit returns a single PullRequest that was merged and introduces the commit on trunk
// trunk is required as multiple PRs can introduce the same commit in the repo, and we need to find the one that did it on the trunk
func (c *GitHubRepoClient) GetPRForCommit(ctx context.Context, commit *gogitplumbing.Hash, trunk string) (*gh.PullRequest, error) {
	prs, _, err := c.PullRequests.ListPullRequestsWithCommit(ctx, c.Owner, c.RepoName, commit.String(), &gh.PullRequestListOptions{})
	if err != nil {
		return nil, err
	}

	var desiredPR *gh.PullRequest = nil
	for _, pr := range prs {
		// TODO for some reason pr.Merged does not return a valid state, so for now we use MergeCommitSHA
		if *pr.Base.Ref == trunk && *pr.MergeCommitSHA != "" {
			desiredPR = pr
			break
		}
	}

	if desiredPR == nil {
		return nil, fmt.Errorf("Could not find any PR that introduces commit %q on the configured trunk %q", commit.String(), trunk)
	}

	return desiredPR, nil
}


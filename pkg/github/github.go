package github

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
	gh "github.com/google/go-github/v44/github"
)

type RepoClient struct {
	Owner    string
	RepoName string
	gh.Client
}

// GetPRForCommit returns a single PullRequest that was merged and introduces the commit on trunk
// trunk is required as multiple PRs can introduce the same commit in the repo, and we need to find the one that did it on the trunk
func (c *RepoClient) GetPRForCommit(ctx context.Context, commit *gogitplumbing.Hash, trunk string) (*gh.PullRequest, error) {
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

func wrapMajor(v *semver.Version) semver.Version {
	return v.IncMajor()
}

func wrapMinor(v *semver.Version) semver.Version {
	return v.IncMinor()
}

func wrapPatch(v *semver.Version) semver.Version {
	return v.IncPatch()
}

func wrapNone(v *semver.Version) semver.Version {
	return *v
}

type SemVerBump string

const (
	Major SemVerBump = "major"
	Minor SemVerBump = "minor"
	Patch SemVerBump = "patch"
	None  SemVerBump = "none"
)

type NoSemVerLabel struct{}

func (*NoSemVerLabel) Error() string {
	return "no GitHub label was found which matches any of the semVer Bump labels"
}

// DetermineSemVerBumpForPR returns a bumping func which can be applied to a semVer Version. It determines the suitable func
// based on a supplied labelmap. This allows users to configure their own labels that the associate with semVer Bumps
func DetermineSemVerBumpForPR(pr *gh.PullRequest, labelMap map[string]SemVerBump) (func(v *semver.Version) semver.Version, error) {
	for _, label := range pr.Labels {

		lm, ok := labelMap[*label.Name]
		// if label does not match supplied map, skip it
		if !ok {
			continue
		}

		switch lm {
		case Major:
			return wrapMajor, nil
		case Minor:
			return wrapMinor, nil
		case Patch:
			return wrapPatch, nil
		case None:
			return wrapNone, nil
		}
	}
	return nil, &NoSemVerLabel{}
}

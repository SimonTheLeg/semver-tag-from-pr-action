package github

import (
	"context"
	"testing"

	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
	gh "github.com/google/go-github/v44/github"
	ghmock "github.com/migueleliasweb/go-github-mock/src/mock"
)

func TestGetPRForCommit2(t *testing.T) {
	trunk := "main"
	mockSha := "e83c5163316f89bfbde7d9ab23ca2e25604af290"
	tt := map[string]struct {
		inPRs   []*gh.PullRequest
		sha     string
		expPRID int64
	}{
		"single PR": {
			inPRs: []*gh.PullRequest{
				{
					ID:             gh.Int64(1),
					Base:           &gh.PullRequestBranch{Ref: &trunk},
					MergeCommitSHA: ptr("asdf1234"),
				},
			},
			sha:     mockSha,
			expPRID: *gh.Int64(1),
		},
		"no PR": {
			inPRs:   []*gh.PullRequest{},
			sha:     mockSha,
			expPRID: *gh.Int64(0),
		},
		"PRs with same sha, but only one is trunk": {
			inPRs: []*gh.PullRequest{
				{
					ID:             gh.Int64(1),
					Base:           &gh.PullRequestBranch{Ref: ptr("not-main-maybe-develop-or-something")},
					MergeCommitSHA: ptr("asdf1234"),
				},
				{
					ID:             gh.Int64(2),
					Base:           &gh.PullRequestBranch{Ref: &trunk},
					MergeCommitSHA: ptr("asdf1234"),
				},
			},
			sha:     mockSha,
			expPRID: *gh.Int64(2),
		},
		"PR with matching sha, but not merged yet": {
			inPRs: []*gh.PullRequest{
				{
					ID:             gh.Int64(1),
					Base:           &gh.PullRequestBranch{Ref: &trunk},
					MergeCommitSHA: ptr(""),
				},
			},
			sha:     mockSha,
			expPRID: *gh.Int64(0),
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {

			mockc := ghmock.NewMockedHTTPClient(
				ghmock.WithRequestMatch(
					ghmock.GetReposCommitsPullsByOwnerByRepoByCommitSha,
					tc.inPRs,
				),
			)
			c := GitHubRepoClient{
				Client:   *gh.NewClient(mockc),
				Owner:    "mock",
				RepoName: "mock",
			}

			hash := gogitplumbing.NewHash(tc.sha)
			pr, err := c.GetPRForCommit(context.Background(), &hash, "main")

			if tc.expPRID != 0 {
				if err != nil {
					t.Errorf("Exp error to be nil, but got %q", err)
				}
				if *pr.ID != tc.expPRID {
					t.Errorf("Exp PRID to be %d, got %d", tc.expPRID, *pr.ID)
				}
			} else {
				if err == nil {
					t.Errorf("Exp PR to not be found, but error is nil")
				}
			}

		})
	}

}

func ptr[T any](v T) *T {
	return &v
}

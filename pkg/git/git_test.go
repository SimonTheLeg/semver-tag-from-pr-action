package git

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
)

func TestGetLatestCommitSha(t *testing.T) {
	// TODO for now hardcode this, until we have found a way to use mocking for this
	branchname := "main"
	repoPath := "/Users/simonbein/github/SimonTheLeg/tag-on-merge-integration-infra"
	expSha := "3084c12d642cf841699ba08e5218102a7450e43b"

	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	res, err := GetCommitForBranch(repo, branchname)
	if err != nil {
		t.Fatal(err)
	}

	if res.Hash().String() != expSha {
		t.Errorf("Exp sha %q, got %q", expSha, res.Hash().String())
	}

}

func TestFindLatestSemVerTag(t *testing.T) {
	repoPath := "/Users/simonbein/github/SimonTheLeg/tag-on-merge-integration-infra"

	repo, err := gogit.PlainOpen(repoPath)
	if err != nil {
		t.Fatal(err)
	}

	latestSemVer, err := FindLatestSemVerTag(repo)
	if err != nil {
		t.Fatal(err)
	}

	expSemver, _ := semver.NewVersion("v1.1.1")

	if !latestSemVer.Equal(expSemver) {
		t.Errorf("Exp semVer %q, got %q", expSemver.String(), latestSemVer.String())
	}

}

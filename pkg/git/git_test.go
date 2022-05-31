package git

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
)

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

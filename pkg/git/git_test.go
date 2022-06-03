package git

import (
	"errors"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

const integrationRepoSha = "0a3f8c254543b1231bde79e5a2483a6b9a3d4081"

func TestFindLatestSemVerTagIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integrationtest TestFindLatestSemVerTagIntegration")
	}

	repo := getIntegrationRepoForRef(t, integrationRepoSha)

	latestSemVer, err := FindLatestSemVerTag(repo)
	if err != nil {
		t.Fatal(err)
	}

	expSemver, _ := semver.NewVersion("v1.1.0")

	if !latestSemVer.Equal(expSemver) {
		t.Errorf("Exp semVer %q, got %q", expSemver.String(), latestSemVer.String())
	}

}

func TestTagExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integrationtest TestTagExists")
	}

	tt := map[string]struct {
		tag string
		exp bool
	}{
		"tag exists": {
			"v1.0.0",
			true,
		},
		"tag does not exists": {
			"does-not-exist",
			false,
		},
	}

	repo := getIntegrationRepoForRef(t, integrationRepoSha)

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			res, err := tagExists(repo, tc.tag)
			if err != nil {
				t.Fatal(err)
			}

			if tc.exp != res {
				t.Errorf("Expected tagExists to be %t, got %t", tc.exp, res)
			}
		})
	}

}

// returns the integration repo without the need to constantly clone it
func getIntegrationRepoForRef(t *testing.T, ref string) *gogit.Repository {
	const repostorespath = "/tmp/tag-on-merge-integration-infra"
	const remotename = "origin"

	// clone the repository, if it is not already present
	var repo *gogit.Repository
	var cloned bool
	if _, err := os.Stat(repostorespath); os.IsNotExist(err) {
		var err error
		repo, err = gogit.PlainClone(repostorespath, false, &gogit.CloneOptions{
			URL: "https://github.com/SimonTheLeg/tag-on-merge-integration-infra.git",
		})
		if err != nil {
			t.Fatal(err)
		}
		cloned = true
	} else {
		repo, err = gogit.PlainOpen(repostorespath)
		if err != nil {
			t.Fatal(err)
		}
	}

	w, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	// check out the desired commit
	err = w.Checkout(&gogit.CheckoutOptions{Hash: plumbing.NewHash(ref)})
	// if the commit does not exists and we have not cloned the repository, pull updates
	if errors.Is(plumbing.ErrReferenceNotFound, err) && cloned == false {
		err = w.Pull(&gogit.PullOptions{RemoteName: remotename})
		if err != nil {
			t.Fatal(err)
		}
		err = w.Checkout(&gogit.CheckoutOptions{Hash: plumbing.NewHash(ref)})
	}
	if err != nil {
		t.Fatal(err)
	}

	return repo
}

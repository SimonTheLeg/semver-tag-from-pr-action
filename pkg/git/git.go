package git

import (
	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
)

// GetCommitForBranch simply returns the Reference a branch points to. It is useful, since most of the
// user inputs work with branchnames, rather than commits
func GetCommitForBranch(repo *gogit.Repository, branchname string) (*gogitplumbing.Reference, error) {

	ref, err := repo.Reference(gogitplumbing.NewBranchReferenceName(branchname), true)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

type NoSemVerTag struct{}

func (e *NoSemVerTag) Error() string { return "no tag matching a semVer version was found" }

// FindLatestSemVerTag returns the latest git Tag which is a semVer version.
// Precedence is defined by semVer as follows: major > minor > patch
func FindLatestSemVerTag(repo *gogit.Repository) (*semver.Version, error) {

	tagrefs, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	startingSemVer, _ := semver.NewVersion("v0.0.0")
	latestSemVer := startingSemVer

	err = tagrefs.ForEach(func(r *gogitplumbing.Reference) error {
		semVer, err := semver.NewVersion(r.Name().Short())
		if err != nil {
			return nil // if tag is not a semVer, skip it
		}

		if semVer.GreaterThan(latestSemVer) {
			latestSemVer = semVer
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if latestSemVer.Equal(startingSemVer) {
		return nil, &NoSemVerTag{}
	}

	return latestSemVer, nil
}

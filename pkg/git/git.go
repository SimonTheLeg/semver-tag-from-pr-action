package git

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
)

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

// SetAnnotatedTag sets an annotated tag in the git repository.
// Tag cannot be empty. If msg is empty, it will be set to the value of tag
func SetAnnotatedTag(repo *gogit.Repository, tag, msg string) error {
	if tag == "" {
		return fmt.Errorf("Tag cannot be empty")
	}
	if msg == "" {
		msg = tag
	}

	tagExists, err := tagExists(repo, tag)
	if err != nil {
		return fmt.Errorf("Could not fetch tag %q: %v", tag, err)
	}
	if tagExists {
		return fmt.Errorf("Tag %q already exists", tag)
	}

	h, err := repo.Head()
	if err != nil {
		return fmt.Errorf("Failed to get HEAD: %v", err)
	}

	_, err = repo.CreateTag(tag, h.Hash(), &gogit.CreateTagOptions{
		Message: msg,
	})
	if err != nil {
		return fmt.Errorf("Could not create tag: %v", err)
	}

	return nil
}

func tagExists(repo *gogit.Repository, tag string) (bool, error) {
	// because prefixes are private in go-git, we have to recreated them https://github.com/go-git/go-git/blob/bf3471db54b0255ab5b159005069f37528a151b7/plumbing/reference.go#L11
	tagPrefix := "refs/tags/"

	tags, err := repo.Tags()
	if err != nil {
		return false, err
	}

	found := false
	err = tags.ForEach(func(r *gogitplumbing.Reference) error {
		if r.Name().String() == tagPrefix+tag {
			found = true
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	return found, nil
}

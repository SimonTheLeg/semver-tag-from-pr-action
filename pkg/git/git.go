package git

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Masterminds/semver/v3"
	gogit "github.com/go-git/go-git/v5"
	gogitconfig "github.com/go-git/go-git/v5/config"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	gogittransport "github.com/go-git/go-git/v5/plumbing/transport"
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
		Tagger: &object.Signature{
			Name:  "semver-tag-from-pr-action",
			Email: "semver-tag-from-pr-action@githubactions.com",
			When:  time.Now(),
		},
		Message: msg,
	})
	if err != nil {
		return fmt.Errorf("Could not create tag: %v", err)
	}

	return nil
}

// PushTag pushes the given tag to the given remote. If remote is empty, 'origin' will be chosen
func PushTag(repo *gogit.Repository, auth gogittransport.AuthMethod, tag string, remote string) error {
	if remote == "" {
		remote = "origin"
	}
	refspec := fmt.Sprintf("refs/tags/%s:refs/tags/%s", tag, tag)

	po := &gogit.PushOptions{
		RemoteName: remote,
		Progress:   os.Stdout,
		RefSpecs:   []gogitconfig.RefSpec{gogitconfig.RefSpec(refspec)},
		Auth:       auth,
	}

	err := repo.Push(po)
	if err != nil {
		if err == gogit.NoErrAlreadyUpToDate {
			log.Println("origin remote was up to date, no tag pushed")
			return nil
		}
		return fmt.Errorf("failed to push tag %q to remote %q: %v", tag, remote, err)
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

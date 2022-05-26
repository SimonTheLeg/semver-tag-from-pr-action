package cmd

import (
	"context"
	"fmt"

	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/config"
	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/git"
	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/github"
)

func Run(conf *config.Config) error {
	sha, err := git.GetCommitForBranch(conf.Repo, conf.Trunk)
	if err != nil {
		return err
	}
	hash := sha.Hash()

	semVerTag, err := git.FindLatestSemVerTag(conf.Repo)
	if err != nil {
		return err
	}

	pr, err := conf.Repoclient.GetPRForCommit(context.Background(), &hash, conf.Trunk)
	if err != nil {
		return err
	}

	bumpFunc, err := github.DetermineSemVerBumpForPR(pr, conf.Labelmap)
	if err != nil {
		return err
	}

	newSemVerTag := bumpFunc(semVerTag)

	fmt.Printf("Old Tag %s, new Tag %s\n", semVerTag.Original(), newSemVerTag.Original())

	return nil
}
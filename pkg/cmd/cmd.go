package cmd

import (
	"context"
	"fmt"

	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/config"
	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/git"
	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/github"
	"github.com/sethvargo/go-githubactions"
)

func Run(conf *config.Config) error {
	semVerTag, err := git.FindLatestSemVerTag(conf.Repo)
	if err != nil {
		return err
	}

	pr, err := conf.Repoclient.GetPRForCommit(context.Background(), conf.EventSha, conf.Trunk)
	if err != nil {
		return err
	}

	bumpFunc, err := github.DetermineSemVerBumpForPR(pr, conf.Labelmap)
	if err != nil {
		return err
	}

	newSemVerTag := bumpFunc(semVerTag)

	fmt.Printf("Old Tag %s, new Tag %s\n", semVerTag.Original(), newSemVerTag.Original())

	githubactions.SetOutput("old-tag", semVerTag.Original())
	githubactions.SetOutput("new-tag", newSemVerTag.Original())

	return nil
}

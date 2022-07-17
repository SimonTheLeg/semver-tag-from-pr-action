package cmd

import (
	"context"
	"fmt"
	"log"

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

	// aka bump none was chosen
	if semVerTag.Equal(&newSemVerTag) {
		log.Println("Detected no bump. Skip setting and pushing new tag")
		return nil
	}

	if !conf.ShouldSetTag {
		log.Println("should_set_tag was set to false. Skip setting and pushing new tag")
		return nil
	}
	err = git.SetAnnotatedTag(conf.Repo, newSemVerTag.Original(), "")
	if err != nil {
		return err
	}

	if !conf.ShouldPushTag {
		log.Println("should_push_tag was set to false. Skip pushing new tag")
		return nil
	}
	err = git.PushTag(conf.Repo, conf.RepoAuth, newSemVerTag.Original(), "")
	if err != nil {
		return err
	}

	return nil
}

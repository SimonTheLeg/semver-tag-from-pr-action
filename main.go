package main

import (
	"fmt"
	"os"

	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/cmd"
	"github.com/SimonTheLeg/semver-tag-on-merge-action/pkg/config"
)

func main() {
	conf, err := config.ConfigInsideActions()
	if err != nil {
		errExit(err)
	}

	err = cmd.Run(conf)
	if err != nil {
		errExit(err)
	}
}

func errExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

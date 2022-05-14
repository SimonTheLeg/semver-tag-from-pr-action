package git

import (
	gogit "github.com/go-git/go-git/v5"
	gogitplumbing "github.com/go-git/go-git/v5/plumbing"
)

func GetCommitForBranch(repo *gogit.Repository, branchname string) (*gogitplumbing.Reference, error) {

	ref, err := repo.Reference(gogitplumbing.NewBranchReferenceName(branchname), true)
	if err != nil {
		return nil, err
	}

	return ref, nil
}

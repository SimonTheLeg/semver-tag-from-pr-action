# Create SemVer Bumps based on PR labels

- [Installation](#installation)
- [Usage](#usage)
- [Customization](#customization)
- [FAQ](#FAQ)
- [Contributing](#contributing)

This action allows you to create new semVer git tags on the trunk branch by using labels on PRs.
The action's main aim is to simplify release processes. Here's an example to illustrate:

1. Current tag on main branch is `v1.0.0`
2. Contributor creates a new PR
3. Maintainer reviews PR and adds the `merge-minor` label (for more options see [Usage section](##Usage))
4. The PR gets merged
5. Action will search through the git history and list of PRs. It determines that the requested semVer bump is minor (second digit) and that the last version was `v1.0.0`
6. Action tags the new commit with the new version `v1.1.0` and pushes the tag

## Installation

The action requires a starting git tag, that is a semVer version. If you are not already using git-tags for your versions, run:

```sh
git tag v1.0.0 # or 1.0.0 or any other valid semVer version
git push origin v1.0.0
```

Minimum Config

```yaml
on:
  push:
    branches:
      - main # can be main or any other release branch
jobs:
  semver_tag_from_pr:
    # permissions required to find PR for last commit and push the new tag
    permissions:
      contents: write
      pull-requests: read
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: "0" # we need full git-history to determine the last semVer tag
      - name: bump semVer
        uses: simontheleg/semver-tag-from-pr-action@v1
        with:
          repo-token: ${{ secrets.GITHUB_TOKEN }}
```

## Usage

After the installation, simply label your PR with one of the following:

* `merge-major` for e.g. v1.0.0 -> v2.0.0
* `merge-minor` for e.g. v1.0.0 -> v1.1.0
* `merge-patch` for e.g. v1.0.0 -> v1.0.1
* `merge-none` no new semVer Tag will be created

Pro tip: The action works well with mergebots when configuring the same labels. Then whenever the label is set, the bot automatically attempts the merge and you will never forget setting the label 🙌

For example for the [kodiak-mergebot](https://kodiakhq.com/) you configure do something like this:

```toml
# kodiak.toml
automerge_label = ["merge-major", "merge-minor", "merge-patch", "merge-none"]
```

## Customization

### Bump Labels

You can specify your own labels, that correspond to the semVer-bumps. For example if you want to adopt Gregor Martynus' `breaking`-`feature`-`fix` notation, you could do:

```yaml
...
with:
  label-major: merge-breaking
  label-minor: merge-feature
  label-patch: merge-fix
  label-none: merge-no-new-version
```

### Disabling Label Set and Push

By default the action will create the new semVer tag and push it into your repository. If instead you want to use your own logic to handle tags:

1. You can disable the setting of the tag inside gh-actions by setting:

    ```yaml
    ...
    with:
      should-set-tag: false
    ```

2. You can disable the push of the tag by setting:

    ```yaml
    ...
    with:
      should-push-tag: false
    ```

In both cases you can use the outputs `old-tag` and `new-tag` of this action in your own jobs:

```yaml
steps:
  ...
  - name: bump semVer
    id: bump-semver # you need to set an id here
    uses: simontheleg/semver-tag-from-pr-action@v1
    with:
      repo-token: ${{ secrets.GITHUB_TOKEN }}
  - name: your step or action
    env: # or alternatively "with", if you are using an action
      old_tag: ${{ steps.bump-semver.outputs.old-tag }}
      new_tag: ${{ steps.bump-semver.outputs.new-tag }}
    run: |
      echo old tag: ${old_tag}
      echo new tag: ${new_tag}
```

## FAQ

**I forgot setting the label on the PR initially and I already merged it**

Assuming you have not merged in any other PRs since then: You can set the label the PR afterwards and then just re-run the action.

## Contributing

### Tests

By default `go test ./...` will run both unit and integration tests. Integration tests clone the [integration-infra repo](https://github.com/SimonTheLeg/semver-tag-from-pr-integration-infra) into `/tmp` once.
As a result, Integration tests should not be run with the `parallel` flag. 

If you only want to run unit-test, you can do so by using the `-short` flag:

```sh
go test -short ./...
```

If you want to only run integration tests, simply run:

```sh
go test -run Integration ./...
```

# Create SemVer Bumps based on PR tags

## Contributing

### Tests

By default `go test ./...` will run both unit and integration tests. Integration tests clone the [integration-infra repo](https://github.com/SimonTheLeg/semver-tag-from-pr-integration-infra) into `/tmp` once.
As a result, Integration tests should not be run with the `parallel` flag. If you only want to run unit-test, you can do so by using the `-short` flag:

```sh
go test -short ./...
```

If you want to only run integration tests, simply run:

```sh
go test -run Integration ./...
```

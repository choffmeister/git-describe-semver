# git-describe-semver

Replacement for `git describe --tags` that produces [semver](https://semver.org/) compatible versions that follow to semver sorting rules.

## Comparison

Corresponding git tag | git describe --tags | git-describe-semver
--- | --- | ---
`v1.2.3` | `v1.2.3` | `v1.2.3`
`v1.2.3` | `v1.2.3-23-gabc1234` | `v1.2.4-dev.23.gabc1234`
`v1.2.3-rc.1` | `v1.2.3-rc.1-23-gabc1234` | `v1.2.3-rc.1.dev.23.gabc1234`

## Usage

```
# binary
# TODO

# with docker
docker run --rm -v $PWD:/workdir choffmeister/git-describe-semver:latest
```

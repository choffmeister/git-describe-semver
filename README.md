# git-describe-semver

Replacement for `git describe --tags` that produces [semver](https://semver.org/) compatible versions that follow to semver sorting rules.

## Comparison

Previous git tag | git describe --tags | git-describe-semver --fallback v0.0.0
--- | --- | ---
`v1.2.3` | `v1.2.3` | `v1.2.3`
`v1.2.3` | `v1.2.3-23-gabc1234` | `v1.2.4-dev.23.gabc1234`
`v1.3.0-rc.1` | `v1.3.0-rc.1-23-gabc1234` | `v1.3.0-rc.1.dev.23.gabc1234`
`v1.3.0-rc.1+info` | `v1.3.0-rc.1+info-23-gabc1234` | `v1.3.0-rc.1.dev.23.gabc1234+info`
none | fail | `v0.0.0-dev.23.gabc1234`

## Usage

* Flag `--fallback v0.0.0`: Fallback to given tag name if no tag is available
* Flag `--drop-prefix`: Drop any present prefix (like `v`) from the output
* Flag `--prerelease-suffix`: Adds a dash-separated suffix to the prerelease part
* Flag `--format`: Changes output (use `<version>` as placeholder)

### Docker

```bash
cd my-git-directory
docker pull ghcr.io/choffmeister/git-describe-semver:latest
docker run --rm -v $PWD:/workdir ghcr.io/choffmeister/git-describe-semver:latest
```

### GitHub action

```yaml
# .github/workflows/build.yml
name: build
jobs:
  update:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
      with:
        fetch-depth: 0
    - uses: docker://ghcr.io/choffmeister/git-describe-semver
      id: git_describe_semver
      with:
        args: --fallback v0.0.0 --drop-prefix --format "::set-output name=version::<version>"
    - run: echo This is the version ${{ steps.git_describe_semver.outputs.version }}
```

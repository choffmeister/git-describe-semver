# git-describe-semver

Replacement for `git describe --tags` that produces [semver](https://semver.org/) compatible versions that follow to semver sorting rules.

## Comparison

Previous git tag | git describe --tags | git-describe-semver --fallback v0.0.0
--- | --- | ---
`v1.2.3` | `v1.2.3` | `v1.2.3`
`v1.2.3` | `v1.2.3-23-gabc1234` | `v1.2.4-dev.23.gabc1234`
`v1.2.3-rc.1` | `v1.2.3-rc.1-23-gabc1234` | `v1.2.3-rc.1.dev.23.gabc1234`
`v1.2.3-rc.1+info` | `v1.2.3-rc.1+info-23-gabc1234` | `v1.2.3-rc.1.dev.23.gabc1234+info`
none | fail | `v0.0.0-dev.23.gabc1234`

## Usage

* Flag `--fallback v0.0.0`: Fallback to given tag name if no tag is available
* Flag `--drop-prefix`: Drop any present prefix (like `v`) from the output

### Binary

```bash
cd my-git-directory
wget -q https://github.com/choffmeister/git-describe-semver/releases/download/v0.2.0/git-describe-semver-linux-amd64
chmod +x git-describe-semver-linux-amd64
./git-describe-semver-linux-amd64 --fallback v0.0.0
```

### Docker

```bash
cd my-git-directory
docker pull choffmeister/git-describe-semver:v0.2.0
docker run --rm -v $PWD:/workdir choffmeister/git-describe-semver:v0.2.0 --fallback v0.0.0
```

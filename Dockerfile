FROM scratch
ENTRYPOINT ["/bin/git-describe-semver"]
COPY git-describe-semver /bin/git-describe-semver
WORKDIR /workdir

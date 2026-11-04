FROM scratch
ARG TARGETPLATFORM
ENTRYPOINT ["/bin/git-describe-semver"]
COPY $TARGETPLATFORM/git-describe-semver /bin/git-describe-semver
WORKDIR /workdir

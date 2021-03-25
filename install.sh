#!/bin/bash
set -eo pipefail

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  VARIANT="linux-amd64"
  mkdir -p "$HOME/bin"
  TARGET="$HOME/bin/git-describe-semver"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  VARIANT="darwin-amd64"
  mkdir -p "$HOME/bin"
  TARGET="$HOME/bin/git-describe-semver"
else
  echo "Unknown OS type $OSTYPE"
  exit 1
fi
LATEST_VERSION="$(curl -s https://api.github.com/repos/choffmeister/git-describe-semver/releases/latest | grep "tag_name" | awk '{print substr($2, 2, length($2)-3)}')"

curl -fsSL -o "$TARGET" "https://github.com/choffmeister/git-describe-semver/releases/download/$LATEST_VERSION/git-describe-semver-$VARIANT"
chmod +x "$TARGET"

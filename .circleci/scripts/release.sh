#!/bin/sh
set -e

GITHUB_NAMESPACE=$1
GITHUB_REPOSITORY=$2
GH_CLI_RELEASES_URL="https://github.com/cli/cli/releases"
FILE_NAME="gh"
BUILD_ARCHITECTURE="linux_amd64.deb"
DELIMETER="_"
PACKAGE_FILE="$FILE_NAME$DELIMETER$BUILD_ARCHITECTURE"

get_release_candidate_tag() {
  git tag --sort=v:refname | grep -E "v[0-9]+\.[0-9]+\.[0-9]+" | tail -n 1
}

RELEASE_CANDIDATE_TAG=$(get_release_candidate_tag)
CURRENT_VERSION="$(printf '%s' "$RELEASE_CANDIDATE_TAG" | cut -c 2-2)"
RELEASE_VERSION=$(if [ "$CURRENT_VERSION" -gt 1 ]; then echo "v$CURRENT_VERSION"; else echo; fi)

release_to_pkg_go() {
  echo "Triggering pkg.go.dev about new release..."
  curl -X POST "https://pkg.go.dev/fetch/github.com/$GITHUB_NAMESPACE/$GITHUB_REPOSITORY/$RELEASE_VERSION@$RELEASE_CANDIDATE_TAG"
}

gh_cli_latest_release() {
  curl -sL -o /dev/null -w '%{url_effective}' "$GH_CLI_RELEASES_URL/latest" | rev | cut -f 1 -d '/'| rev
}

download_gh_cli() {
  test -z "$VERSION" && VERSION="$(gh_cli_latest_release)"
  test -z "$VERSION" && {
    echo "Unable to get GitHub CLI release." >&2
    exit 1
  }
  curl -s -L -o "$PACKAGE_FILE" "$GH_CLI_RELEASES_URL/download/$VERSION/$FILE_NAME$DELIMETER$(printf '%s' "$VERSION" | cut -c 2-100)$DELIMETER$BUILD_ARCHITECTURE"
}

install_gh_cli() {
  sudo dpkg -i "$PACKAGE_FILE"
  rm "$PACKAGE_FILE"
}

release_to_github() {
  echo "Downloading and installing latest gh cli..."
  download_gh_cli
  install_gh_cli
  echo "Publishing new release notes to GitHub..."
  gh release create "$RELEASE_CANDIDATE_TAG" --generate-notes
}

release_to_pkg_go
release_to_github

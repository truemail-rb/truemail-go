#!/bin/sh
set -e

SEMVER_REGEX_PATTERN="[0-9]+\.[0-9]+\.[0-9]+"

latest_changelog_tag() {
  grep -Po "(?<=\#\# \[)$SEMVER_REGEX_PATTERN?(?=\])" CHANGELOG.md | head -n 1
}

latest_git_tag() {
  git tag --sort=v:refname | grep -E "v$SEMVER_REGEX_PATTERN" | tail -n 1
}

TAG_CANDIDATE="v$(latest_changelog_tag)"

if [ "$TAG_CANDIDATE" != "$(latest_git_tag)" ]
then
  echo "Configuring git..."
  git config --global user.email "${PUBLISHER_EMAIL}"
  git config --global user.name "${PUBLISHER_NAME}"
  echo "Pushing new semver tag to GitHub..."
  git tag "$TAG_CANDIDATE"
  git push --tags
  echo "Updating develop branch with new semver tag..."
  git checkout develop
  git merge "$TAG_CANDIDATE" --ff --no-edit
  git push origin develop
else
  echo "Latest changelog tag ($TAG_CANDIDATE) already released on GitHub. Tagging is not required."
fi

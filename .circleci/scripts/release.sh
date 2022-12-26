#!/bin/sh
set -e

latest_tag() {
  git tag -l | grep -E "^v[0-9]+\.[0-9]+\.[0-9]+" | cut -d"-" -f 1 | sort | tail -n 1
}

publish_release() {
  echo "Triggering pkg.go.dev about new truemail-go release..."
  curl -X POST "https://pkg.go.dev/fetch/github.com/truemail-rb/truemail-go@$(latest_tag)"
}

publish_release

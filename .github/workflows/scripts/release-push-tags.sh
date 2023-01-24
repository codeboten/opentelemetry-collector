#!/bin/bash -ex

git config user.name "$GITHUB_ACTOR"
git config user.email "$GITHUB_ACTOR@users.noreply.github.com"

make push-tags MODSET="${MODSET}"

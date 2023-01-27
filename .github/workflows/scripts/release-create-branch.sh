#!/bin/bash -ex

git config user.name "$GITHUB_ACTOR"
git config user.email "$GITHUB_ACTOR@users.noreply.github.com"

BRANCH="release/v${CANDIDATE_STABLE}-v${CANDIDATE_BETA}"
git checkout -b "${BRANCH}"
git push origin "${BRANCH}"

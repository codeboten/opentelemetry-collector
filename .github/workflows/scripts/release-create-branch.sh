#!/bin/bash -ex

git config user.name opentelemetrybot
git config user.email 107717825+opentelemetrybot@users.noreply.github.com

git checkout "${COMMIT}"
BRANCH="release/v${CANDIDATE_STABLE}-v${CANDIDATE_BETA}"
git checkout -b "${BRANCH}"
git push "git@github.com:${REPO}" "${BRANCH}"

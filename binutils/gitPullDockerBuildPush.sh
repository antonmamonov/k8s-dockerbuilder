#!/bin/bash

# check if GIT_PULL_REPO_URL is set
if [ -z "$GIT_PULL_REPO_URL" ]; then
    echo "GIT_PULL_REPO_URL is not set"
    exit 1
fi

# check if DOCKER_IMAGE_DESTINATION is set
if [ -z "$DOCKER_IMAGE_DESTINATION" ]; then
    echo "DOCKER_IMAGE_DESTINATION is not set (ie, hello/world:v0.0.1)"
    exit 1
fi

# check if DOCKER_CONFIG_JSON is set
if [ -z "$DOCKER_CONFIG_JSON" ]; then
    echo "DOCKER_CONFIG_JSON is not set"
    exit 1
fi

cd /workspace

# remove any previous git repos
rm -fr /workspace/apptobuild
git clone $GIT_PULL_REPO_URL /workspace/apptobuild

# if GIT_BRANCH exists then checkout that branch
if [ -n "$GIT_BRANCH" ]; then
    cd /workspace/apptobuild
    git checkout $GIT_BRANCH
fi

# save the docker config json to tmp
mkdir -p /root/.docker
rm -fr /root/.docker/config.json
echo $DOCKER_CONFIG_JSON > /root/.docker/config.json

cd /workspace/apptobuild
kaniko --context /workspace/apptobuild \
    --single-snapshot \
    --push-retry 7 \
    --destination $DOCKER_IMAGE_DESTINATION
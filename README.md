# Build & push Docker images inside your Kubernetes cluster

<img src='./images/containers.png' width="420" height="420">

## Prerequisite

Given the CLI needs access to your Kubernetes cluster to build an image it would require an authenticated context. The CLI tool assumes a working Kubernetes config file at `$HOME/.kube/config`.

## Required Environment Variables

- `GIT_PULL_REPO_URL` - The URL of the Git Repository to pull (ie, https://username:token@github.com/yourorg/repo.git)
- `DOCKER_IMAGE_DESTINATION` - The Docker Image Destination (ie, yourorg/repo:v0.0.1) to push after building
- `DOCKER_CONFIG_JSON` - The Docker config JSON file (ie, `cat ~/.docker/config.json`) to use for authentication

## Optional Environment Variables

- `GIT_BRANCH` - The Git Branch to checkout and build from. By default it assumes the default branch on your external git repository

## Quick start with a simple docker build locally

```bash
# build docker image
./buildDockerImage.sh

# run docker container
./runDockerDev.sh

# inside the docker container
./binutils/sampleKanikoUbuntuDockerBuild.sh
```
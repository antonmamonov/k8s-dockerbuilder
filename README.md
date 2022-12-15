# Build & Push Docker images inside your Kubernetes cluster

<img src='./images/containers.png' width="420" height="420">

## Prerequisites for kubebuild CLI tool

1. Given the CLI needs access to your Kubernetes cluster to build an image it would require an authenticated context. The CLI tool assumes a working Kubernetes config file at `$HOME/.kube/config`.

2. To push the built docker image to your private container registry the CLI tool assumes a valid docker config json at `$HOME/.docker/config`

## Quick run kubebuild CLI

```bash
# build the cli for your environment if needed
go build -o kubebuild .

# copy over the executable to your system's bin folder
sudo cp ./kubebuild /usr/local/bin/kubebuild

# See all commands
kubebuild help

# navigate to any local git repo and give the destination docker image tag as an argument (ie, hello/image:v0.0.1)
kubebuild build hello/image:v0.0.1

# If successful you should see a 'dockerbuild' job in your cluster
kubectl get jobs 
```

## Required Environment Variables for Docker Container

- `GIT_PULL_REPO_URL` - The URL of the Git Repository to pull (ie, https://username:token@github.com/yourorg/repo.git)
- `DOCKER_IMAGE_DESTINATION` - The Docker Image Destination (ie, yourorg/repo:v0.0.1) to push after building
- `DOCKER_CONFIG_JSON` - The Docker config JSON file (ie, `cat ~/.docker/config.json`) to use for authentication

## Optional Environment Variables for Docker Container

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
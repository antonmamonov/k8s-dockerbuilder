#!/bin/bash

kaniko --dockerfile /app/sampledockerfiles/ubuntu.Dockerfile --destination antonm/sampleubuntutext:v0.0.1 --no-push
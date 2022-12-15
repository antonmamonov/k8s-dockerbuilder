#!/bin/bash

IMAGENAME=antonm/dockerbuilder:v0.0.1
CONTAINERNAME=dockerbuilder

docker rm -f $CONTAINERNAME

docker run -t -d \
    --name $CONTAINERNAME \
    -v $PWD:/app \
    $IMAGENAME

docker exec -it $CONTAINERNAME /bin/bash
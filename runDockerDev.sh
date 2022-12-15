#!/bin/bash

IMAGENAME=antonm/dockerbuilder
CONTAINERNAME=dockerbuilder

docker rm -f $CONTAINERNAME

docker run -t -d \
    --name $CONTAINERNAME \
    -v $PWD:/app \
    $IMAGENAME

docker exec -it $CONTAINERNAME /bin/bash
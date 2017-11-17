#!/bin/bash 

source build/build.env.sh
$DOCKER login -u="$DOCKER_USERNAME" -p="$DOCKER_PASSWORD"
$DOCKER push $DOCKER_IMAGE


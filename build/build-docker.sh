#!/bin/bash -x

source build/build.env.sh


$DOCKER build --no-cache -t $DOCKER_IMAGE . 


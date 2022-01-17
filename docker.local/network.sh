#!/bin/bash

docker network create \
--driver=bridge \
--subnet=172.16.238.0/24 \
sakura-network

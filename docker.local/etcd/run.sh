#!/bin/bash

cd ./docker.local/etcd

rm -rfv ./data

docker-compose -p etcd-cluster up

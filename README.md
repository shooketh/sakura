# Sakura
Sakura is a network service for generating unique ID numbers inspired by Twitter's Snowflake.

## How to run sakura-cluster and etcd-cluster on local
- install docker https://docs.docker.com/engine/install/
- run `./build/build.sh` to build `sakura` image on local
- run `./docker.local/ipconfig.sh` (macOS)
- run `./docker.local/network.sh` to create docker network `sakura-network`
- run `./docker.local/etcd/run.sh` to run `etcd-cluster`
- run `./docker.local/sakura/run.sh` to run `sakura-cluster`

## How to run client to test it on local
- `cd ./cmd/client`
- `go run main.go`

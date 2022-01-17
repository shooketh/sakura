#!/bin/bash
set -e

build="build"

for arg in "$@"; do
  case $arg in
    -m1|--m1|m1)
      echo "The build will be performed for Apple M1 chip"
      build="buildx build --platform linux/amd64"
      shift
      ;;
  esac
done

docker "${build}" -t sakura -f ./build/dockerfile .

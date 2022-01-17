#!/bin/bash

# macOS platform
if [ "$(uname)" == "Darwin" ]
then
  # etcd
  sudo ifconfig lo0 alias 172.16.238.101
  sudo ifconfig lo0 alias 172.16.238.102
  sudo ifconfig lo0 alias 172.16.238.103
  # sakura
  sudo ifconfig lo0 alias 172.16.238.11
  sudo ifconfig lo0 alias 172.16.238.12
  sudo ifconfig lo0 alias 172.16.238.13
  sudo ifconfig lo0 alias 172.16.238.14
  sudo ifconfig lo0 alias 172.16.238.15
  sudo ifconfig lo0 alias 172.16.238.16
  sudo ifconfig lo0 alias 172.16.238.17
  sudo ifconfig lo0 alias 172.16.238.18
  sudo ifconfig lo0 alias 172.16.238.19
  sudo ifconfig lo0 alias 172.16.238.20
  sudo ifconfig lo0 alias 172.16.238.21
  sudo ifconfig lo0 alias 172.16.238.22
  sudo ifconfig lo0 alias 172.16.238.23
  sudo ifconfig lo0 alias 172.16.238.24
  sudo ifconfig lo0 alias 172.16.238.25
  sudo ifconfig lo0 alias 172.16.238.26
  sudo ifconfig lo0 alias 172.16.238.27
  sudo ifconfig lo0 alias 172.16.238.28
  sudo ifconfig lo0 alias 172.16.238.29
  sudo ifconfig lo0 alias 172.16.238.30
  sudo ifconfig lo0 alias 172.16.238.31
  sudo ifconfig lo0 alias 172.16.238.32
  sudo ifconfig lo0 alias 172.16.238.33
  sudo ifconfig lo0 alias 172.16.238.34
  sudo ifconfig lo0 alias 172.16.238.35
  sudo ifconfig lo0 alias 172.16.238.36
  sudo ifconfig lo0 alias 172.16.238.37
  sudo ifconfig lo0 alias 172.16.238.38
  sudo ifconfig lo0 alias 172.16.238.39
  sudo ifconfig lo0 alias 172.16.238.40
  sudo ifconfig lo0 alias 172.16.238.41
  sudo ifconfig lo0 alias 172.16.238.42
  sudo ifconfig lo0 alias 172.16.238.43
fi

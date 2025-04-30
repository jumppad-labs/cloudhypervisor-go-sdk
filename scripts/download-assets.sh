#!/bin/bash
set -ex

mkdir -p examples/files

if [ ! -f examples/files/ubuntu.img ]; then
  wget -O examples/files/ubuntu.img http://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img
fi

if [ ! -f examples/files/vmlinux ]; then
  wget -O examples/files/vmlinux https://github.com/cloud-hypervisor/linux/releases/download/ch-release-v6.2-20240908/vmlinux
fi
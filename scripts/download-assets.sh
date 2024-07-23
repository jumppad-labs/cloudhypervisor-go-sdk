#!/bin/bash
set -x

if [ ! -f examples/files/noble.img ]; then
  curl -L -o examples/files/noble.img https://cloud-images.ubuntu.com/releases/noble/release/ubuntu-24.04-server-cloudimg-amd64.img
fi

if [ ! -f examples/files/initrd ]; then
  curl -L -o examples/files/initrd.img https://cloud-images.ubuntu.com/releases/noble/release/unpacked/ubuntu-24.04-server-cloudimg-amd64-initrd-generic
fi

if [ ! -f examples/files/vmlinuz ]; then
  curl -L -o examples/files/vmlinuz https://cloud-images.ubuntu.com/releases/noble/release/unpacked/ubuntu-24.04-server-cloudimg-amd64-vmlinuz-generic
fi
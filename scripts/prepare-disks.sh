#!/bin/bash
mkdir -p examples/files
cd examples/files

qemu-img create -f qcow2 -F qcow2 -b ubuntu.img root.img 10G
qemu-img convert -p -f qcow2 -O raw root.img root.raw
#!/bin/bash
set -x

for IMG in $(find examples/files -name '*.img'); do 
  NAME=$(echo $IMG | cut -d '-' -f 1)
  if [ ! -f $NAME.raw ]; then
    qemu-img convert -p -f qcow2 -O raw $IMG $NAME.raw
  fi
done
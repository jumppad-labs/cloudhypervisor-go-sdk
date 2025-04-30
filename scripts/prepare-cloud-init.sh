#!/bin/bash
rm /tmp/cloudinit.img || true

mkdir -p examples/files/cloud-init
cd examples/files/cloud-init

mkdosfs -n CIDATA -C /tmp/cloudinit.img 8192
mcopy -oi /tmp/cloudinit.img -s user-data ::
mcopy -oi /tmp/cloudinit.img -s meta-data ::
mcopy -oi /tmp/cloudinit.img -s network-config ::
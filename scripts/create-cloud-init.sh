#!/bin/bash
set -x

rm -f /tmp/cloudinit.iso
mkdosfs -n CIDATA -C /tmp/cloudinit.iso 8192
mcopy -oi /tmp/cloudinit.iso -s cloud-init/user-data ::
mcopy -oi /tmp/cloudinit.iso -s cloud-init/meta-data ::
mcopy -oi /tmp/cloudinit.iso -s cloud-init/network-config ::
package main

import (
	"context"
	"log"
	"path/filepath"

	sdk "github.com/jumppad-labs/cloudhypervisor-go-sdk"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/api"
)

func main() {
	ctx := context.Background()

	logger := log.Default()

	mac := "12:34:56:78:90:01"
	tap := "tap0"
	ip := "192.168.249.1"
	mask := "255.255.255.0"

	kernel, err := filepath.Abs("examples/files/vmlinux")
	if err != nil {
		logger.Fatal(err)
	}

	disk, err := filepath.Abs("examples/files/root.raw")
	if err != nil {
		logger.Fatal(err)
	}

	cloudinit := "/tmp/cloudinit.img"

	args := "console=ttyS0 ds=nocloud root=/dev/vda1 rw"

	config := api.VmConfig{
		Payload: api.PayloadConfig{
			Kernel:  &kernel,
			Cmdline: &args,
		},
		Disks: &[]api.DiskConfig{
			{
				Path: &disk,
			},
			{
				Path: &cloudinit,
			},
		},
		Net: &[]api.NetConfig{
			{
				Tap:  &tap,
				Mac:  &mac,
				Mask: &mask,
				Ip:   &ip,
			},
		},
		Cpus: &api.CpusConfig{
			BootVcpus: 1,
			MaxVcpus:  1,
		},
		Memory: &api.MemoryConfig{
			Size: 1024 * 1000 * 1000, // 1GB
		},
		Serial: &api.ConsoleConfig{
			Mode: "Tty",
		},
	}

	machine, err := sdk.NewMachine(ctx, config, logger)
	if err != nil {
		logger.Fatal(err)
	}

	err = machine.Start(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	err = machine.Wait(ctx)
	if err != nil {
		logger.Fatal(err)
	}
}

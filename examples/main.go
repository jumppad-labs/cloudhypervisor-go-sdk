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

	username := "erik"
	password := "$6$7125787751a8d18a$sHwGySomUA1PawiNFWVCKYQN.Ec.Wzz0JtPPL1MvzFrkwmop2dq7.4CYf03A5oemPQ4pOFCCrtCelvFBEle/K." // cloud123

	gateway := "192.168.249.1"
	cidr := "192.168.249.2/24"
	mac := "12:34:56:78:90:01"

	// use this firmware if no kernel is specified
	kernel, err := filepath.Abs("examples/files/vmlinuz")
	if err != nil {
		logger.Fatal(err)
	}

	initrd, err := filepath.Abs("examples/files/initrd")
	if err != nil {
		log.Fatal(err)
	}

	disk, err := filepath.Abs("examples/files/noble.raw")
	if err != nil {
		logger.Fatal(err)
	}

	cloudinit, err := sdk.CreateCloudInitDisk("microvm", mac, cidr, gateway, username, password)
	if err != nil {
		logger.Fatal(err)
	}

	_ = cloudinit

	args := "root=/dev/vda1 ro console=tty1 console=ttyS0"
	serial := "/tmp/serial"

	// readonly := true

	config := api.VmConfig{
		Payload: api.PayloadConfig{
			Kernel:    &kernel,
			Initramfs: &initrd,
			Cmdline:   &args,
		},
		Disks: &[]api.DiskConfig{
			{
				Path: disk,
				// Readonly: &readonly,
			},
			// {
			// 	Path: cloudinit,
			// },
		},
		Net: &[]api.NetConfig{
			{
				Mac: &mac,
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
			Mode: "File",
			File: &serial,
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

package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"

	sdk "github.com/jumppad-labs/cloudhypervisor-go-sdk"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"
	"github.com/kr/pretty"
)

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	username := "jumppad"
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

	args := "root=/dev/vda1 ro console=tty1 console=ttyS0"
	serial := "/dev/serial"

	config := client.VmConfig{
		Payload: client.PayloadConfig{
			Kernel:    &kernel,
			Initramfs: &initrd,
			Cmdline:   &args,
		},
		Disks: &[]client.DiskConfig{
			{
				Path: disk,
			},
			{
				Path: cloudinit,
			},
		},
		Net: &[]client.NetConfig{
			{
				Mac: &mac,
			},
		},
		Cpus: &client.CpusConfig{
			BootVcpus: 1,
			MaxVcpus:  1,
		},
		Memory: &client.MemoryConfig{
			Size: 1024 * 1000 * 1000, // 1GB
		},
		Serial: &client.ConsoleConfig{
			Mode: "File",
			File: &serial,
		},
	}

	pretty.Println(config)

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

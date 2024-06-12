package main

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml https://raw.githubusercontent.com/cloud-hypervisor/cloud-hypervisor/master/vmm/src/api/openapi/cloud-hypervisor.yaml

import (
	"context"
	"log"
	"path/filepath"

	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"
)

func main() {
	ip := "192.168.10.10"
	mask := "255.255.255.0"
	mac := "12:34:56:78:90:01"
	tap := "tap0"

	kernel, err := filepath.Abs("examples/hypervisor-fw")
	if err != nil {
		panic(err)
	}

	disk, err := filepath.Abs("examples/focal-server-cloudimg-amd64.raw")
	if err != nil {
		panic(err)
	}

	config := client.VmConfig{
		Payload: client.PayloadConfig{
			Kernel: &kernel,
		},
		Disks: &[]client.DiskConfig{
			{
				Path: disk,
			},
		},
		Net: &[]client.NetConfig{
			{
				Ip:   &ip,
				Mask: &mask,
				Mac:  &mac,
				Tap:  &tap,
			},
		},
		Cpus: &client.CpusConfig{
			BootVcpus: 1,
			MaxVcpus:  1,
		},
		Memory: &client.MemoryConfig{
			Size: 1024 * 1000 * 1000,
		},
	}

	ctx := context.Background()
	machine, err := NewMachine(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	err = machine.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = machine.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

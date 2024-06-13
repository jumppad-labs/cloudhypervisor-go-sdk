package main

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=client/config.yaml https://raw.githubusercontent.com/cloud-hypervisor/cloud-hypervisor/master/vmm/src/api/openapi/cloud-hypervisor.yaml

import (
	"context"
	"fmt"
	"log"
	"net"
	"path/filepath"

	sdk "github.com/jumppad-labs/cloudhypervisor-go-sdk"
)

func main() {
	ctx := context.Background()

	username := "instruqt"
	password := "$6$2XC6sDcIdykdJMyp$j0IIMBPLavRisH.bkFbetP18R.a4IyKctUZ6.84Qw/6ADUMQ074Dp01VZIbYVPwe7SmaPEWmuQKM2UCp.I2At."

	/*
		cloud123 (512) = $6$SR9/pN.80DvU7P97$ap6rBBaN6GdDaQQUOivGzTahjnANXW6Yzwsu42Eit4GrGResGXbuI28a7rge4G3Qug7NKqujFRWGPHOuKe0cl/
		cloud123 (???) = $6$7125787751a8d18a$sHwGySomUA1PawiNFWVCKYQN.Ec.Wzz0JtPPL1MvzFrkwmop2dq7.4CYf03A5oemPQ4pOFCCrtCelvFBEle/K.
	*/

	gateway := "10.0.5.1"
	cidr := "10.0.5.0/24"
	mac := "12:34:56:78:90:01"
	tap := "tap0"

	address, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	ip := address.String()
	mask := network.Mask.String()

	// use this firmware if no kernel is specified
	kernel, err := filepath.Abs("examples/files/hypervisor-fw")
	if err != nil {
		panic(err)
	}

	disk, err := filepath.Abs("examples/files/focal-server-cloudimg-amd64.raw")
	if err != nil {
		panic(err)
	}

	userdata := fmt.Sprintf(`#cloud-config
	users:
		- name: %s
			passwd: %s
			sudo: ALL=(ALL) NOPASSWD:ALL
			lock_passwd: False
			inactive: False
			shell: /bin/bash

	ssh_pwauth: True
	`, username, password)

	err = sdk.CreateCloudInitDisk("microvm", "test", mac, cidr, gateway, userdata)
	if err != nil {
		panic(err)
	}

	config := sdk.Config{
		Kernel: sdk.KernelConfig{
			Path: &kernel,
		},
		Disks: []sdk.DiskConfig{
			{
				Path: disk,
			},
		},
		Network: []sdk.NetworkConfig{
			{
				IpAddress:     &ip,
				MacAddress:    &mac,
				Mask:          &mask,
				HostInterface: &tap,
				Gateway:       &gateway,
			},
		},
		CPU: sdk.CPUConfig{
			BootVcpus: 1,
			MaxVcpus:  1,
		},
		Memory: sdk.MemoryConfig{
			Size: 1024 * 1000 * 1000,
		},
	}

	machine, err := sdk.NewMachine(ctx, config)
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

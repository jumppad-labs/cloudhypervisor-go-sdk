package main

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml https://raw.githubusercontent.com/cloud-hypervisor/cloud-hypervisor/master/vmm/src/api/openapi/cloud-hypervisor.yaml

import (
	"path/filepath"

	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/types"

	"github.com/kr/pretty"
)

func main() {
	vm, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	ip := "192.168.10.10"
	mask := "255.255.255.0"
	mac := "12:34:56:78:90:01"

	kernel, err := filepath.Abs("examples/hypervisor-fw")
	if err != nil {
		panic(err)
	}

	disk, err := filepath.Abs("examples/focal-server-cloudimg-amd64.raw")
	if err != nil {
		panic(err)
	}

	config := types.Config{
		Kernel: types.KernelConfig{
			Path: &kernel,
		},
		Disks: []types.DiskConfig{
			{
				Path: disk,
			},
		},
		Network: []types.NetworkConfig{
			{
				IP:   &ip,
				Mask: &mask,
				MAC:  &mac,
			},
		},
		CPU: types.CPUConfig{
			BootVcpus: 1,
			MaxVcpus:  1,
		},
	}

	info, err := vm.Create(config)
	if err != nil {
		panic(err)
	}

	pretty.Println(info)

	vm.Boot()

	err = vm.Ping()
	if err != nil {
		panic(err)
	}
}

package main

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml https://raw.githubusercontent.com/cloud-hypervisor/cloud-hypervisor/master/vmm/src/api/openapi/cloud-hypervisor.yaml

import (
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/api"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"

	"github.com/kr/pretty"
)

func main() {
	vm, err := client.NewClient()
	if err != nil {
		panic(err)
	}

	info, err := vm.Create(api.VmConfig{})
	if err != nil {
		panic(err)
	}

	pretty.Println(info)

	vm.Boot()

	ping, err := vm.Ping()
	if err != nil {
		panic(err)
	}

	pretty.Println(ping)
}

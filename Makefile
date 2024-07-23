generate:
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	go generate -v
	go mod tidy
	
build:
	go build -o bin/cloudhypervisor-go-sdk examples/main.go

run:
	sudo $(PWD)/bin/cloudhypervisor-go-sdk

kill:
	sudo rm -rf /tmp/cloudinit* || true
	sudo rm /dev/serial || true
	sudo killall cloud-hypervisor || true
	sudo killall cloudhypervisor-go-sdk || true

assets:
	sudo scripts/download-assets.sh
	sudo scripts/create-raw-disks.sh
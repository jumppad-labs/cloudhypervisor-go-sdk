generate:
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	go generate -v
	go mod tidy
	
build:
	go build -o bin/cloudhypervisor-go-sdk examples/main.go

run: build disk cloudinit
	$(PWD)/bin/cloudhypervisor-go-sdk

clean:
	sudo rm examples/files/*.raw || true
	sudo rm /tmp/cloudinit.img || true
	sudo rm -rf /tmp/cloudinit* || true

kill:
	sudo killall cloud-hypervisor || true
	sudo killall cloudhypervisor-go-sdk || true

download:
	scripts/download-assets.sh

disk:
	scripts/prepare-disks.sh

cloudinit:
	scripts/prepare-cloud-init.sh

assets: download disk cloudinit
	
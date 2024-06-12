generate:
	go get github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen
	go generate -v
	go mod tidy
	
build:
	go build -o bin/

run:
	sudo $(PWD)/bin/cloudhypervisor-go-sdk

kill:
	sudo killall cloud-hypervisor
generate:
	go generate -v
build:
	go build -o bin/

run:
	sudo $(PWD)/bin/cloudhypervisor-go-sdk

kill:
	sudo killall cloud-hypervisor
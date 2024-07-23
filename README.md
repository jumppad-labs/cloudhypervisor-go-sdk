# Cloud Hypervisor Golang SDK

## Example

Download and prepare assets.

```shell
make assets
```

Build the example code.

```shell
make build
```

Run the example code.

```shell
make run
```

Confirm that the virtual machine is running.

```shell
sudo curl -s --unix-socket /tmp/cloud-hypervisor.sock \
  http://localhost/api/v1/vm.info | jq .
```

Cleanup.

```shell
make kill
```

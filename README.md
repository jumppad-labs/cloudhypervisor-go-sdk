# Cloud Hypervisor Golang SDK

## Example

Setup cloud-hypervisor.

```shell
sudo setcap cap_net_admin+ep $(which cloud-hypervisor)
```

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

Login using `jumppad` as user and `cloud123` as password.

```shell
Ubuntu 24.04 LTS microvm hvc0

microvm login:
```

From another terminal session, confirm that the virtual machine is running.

```shell
sudo curl -s --unix-socket /tmp/cloud-hypervisor.sock \
  http://localhost/api/v1/vm.info | jq .
```

To exit and cleanup. Either shutdown the machine gracefully or forcefully.

```shell
# Gracefully from inside the vm.
jumppad@microvm:~$ sudo shutdown -h now

# Then after a few seconds, exit the process with CTRL+C
```

```shell
# Forcefully from outside the vm.
make kill
```

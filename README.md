```shell
curl -s --unix-socket /tmp/cloud-hypervisor.sock -i \
  -X PUT 'http://localhost/api/v1/vm.create'  \
  -H 'Accept: application/json'               \
  -H 'Content-Type: application/json'         \
  -d @examples/vmconfig.json
```

```shell
sudo curl -s --unix-socket /tmp/cloud-hypervisor.sock \
  http://localhost/api/v1/vm.info | jq .
```
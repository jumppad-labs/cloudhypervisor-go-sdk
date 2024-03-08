package main

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml https://raw.githubusercontent.com/cloud-hypervisor/cloud-hypervisor/master/vmm/src/api/openapi/cloud-hypervisor.yaml

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/instruqt/cloudhypervisor-go-sdk/api"
	"github.com/kr/pretty"
)

func main() {
	const socket string = "/tmp/cloud-hypervisor.sock"
	const url string = "http://localhost/api/v1/"

	unixClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socket)
			},
		},
	}

	client, err := api.NewClient(url, api.WithHTTPClient(unixClient))
	if err != nil {
		panic(err)
	}

	// client.CreateVM(ctx context.Context, body CreateVMJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.BootVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.PauseVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.PutVmSnapshot(ctx context.Context, body PutVmSnapshotJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.ResumeVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// client.CreateVM(ctx context.Context, body CreateVMJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.PutVmRestore(ctx context.Context, body PutVmRestoreJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// client.GetVmInfo(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// client.PowerButtonVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.ShutdownVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
	// client.DeleteVM(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	ping, err := client.GetVmmPing(context.Background())
	if err != nil {
		panic(err)
	}

	data, err := io.ReadAll(ping.Body)
	if err != nil {
		panic(err)
	}

	pretty.Println(string(data))
}

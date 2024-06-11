package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"

	"github.com/jumppad-labs/cloudhypervisor-go-sdk/api"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/types"
)

type Option func(*ClientImpl) error

const (
	defaultSocket = "/tmp/cloud-hypervisor.sock"
	defaultURL    = "http://localhost/api/v1/"
)

type Client interface {
	Create(config types.Config) (*types.VM, error)
	Boot() (*types.VM, error)
	Pause() (*types.VM, error)
	Resume() (*types.VM, error)
	Snapshot(destination string) error
	Restore(source string) error
	Reboot() (*types.VM, error)
	PowerButton() (*types.VM, error)
	Shutdown() (*types.VM, error)
	Delete() error
	Info() (*types.VM, error)
	Ping() error
}

type ClientImpl struct {
	context    context.Context
	url        string
	socket     string
	unixClient *http.Client
	apiClient  *api.Client
}

func NewClient() (Client, error) {
	path, err := exec.LookPath("cloud-hypervisor")
	if err != nil {
		return nil, err
	}

	go func() {
		cmd := exec.Command(path, "--api-socket", defaultSocket)
		stdErr, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(string(stdErr))
			log.Fatal(err)
		}

		waitErr := cmd.Wait()
		if err != nil {
			log.Fatal(waitErr)
		}
	}()

	unixClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", defaultSocket)
			},
		},
	}

	apiClient, err := api.NewClient(defaultURL, api.WithHTTPClient(unixClient))
	if err != nil {
		return nil, err
	}

	return &ClientImpl{
		context:    context.Background(),
		url:        defaultURL,
		socket:     defaultSocket,
		unixClient: unixClient,
		apiClient:  apiClient,
	}, nil
}

func WithURL(url string) Option {
	return func(c *ClientImpl) error {
		apiClient, err := api.NewClient(url, api.WithHTTPClient(c.unixClient))
		if err != nil {
			return err
		}

		c.apiClient = apiClient
		return nil
	}
}

func WithSocket(socket string) Option {
	return func(c *ClientImpl) error {
		unixClient := &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial("unix", socket)
				},
			},
		}

		apiClient, err := api.NewClient(c.url, api.WithHTTPClient(unixClient))
		if err != nil {
			return err
		}

		c.apiClient = apiClient
		return nil
	}
}

func WithClient(client *http.Client) Option {
	return func(c *ClientImpl) error {
		apiClient, err := api.NewClient(c.url, api.WithHTTPClient(client))
		if err != nil {
			return err
		}

		c.apiClient = apiClient
		return nil
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *ClientImpl) error {
		c.context = ctx
		return nil
	}
}

func (c *ClientImpl) Create(config types.Config) (*types.VM, error) {

	cfg := api.VmConfig{}
	resp, err := c.apiClient.CreateVM(c.context, cfg)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("could not create vm: %s", string(body))
	}

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Boot() (*types.VM, error) {
	_, err := c.apiClient.BootVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Pause() (*types.VM, error) {
	_, err := c.apiClient.PauseVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404, 405

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Resume() (*types.VM, error) {
	_, err := c.apiClient.ResumeVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404, 405

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Snapshot(destination string) error {
	config := api.VmSnapshotConfig{}
	_, err := c.apiClient.PutVmSnapshot(c.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) Restore(source string) error {
	config := api.RestoreConfig{}
	_, err := c.apiClient.PutVmRestore(c.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) Reboot() (*types.VM, error) {
	_, err := c.apiClient.RebootVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404, 405

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) PowerButton() (*types.VM, error) {
	_, err := c.apiClient.PowerButtonVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404, 405

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Shutdown() (*types.VM, error) {
	_, err := c.apiClient.ShutdownVM(c.context)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204, 404, 405

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Delete() error {
	resp, err := c.apiClient.DeleteVM(c.context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not delete vm: %s", string(body))
	}

	return nil
}

func (c *ClientImpl) Info() (*types.VM, error) {
	resp, err := c.apiClient.GetVmInfo(c.context)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("could not get vm info: %s", string(body))
	}

	info := api.VmInfo{}
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	config := types.VmConfigToConfig(info.Config)

	vm := &types.VM{
		Config: config,
	}

	return vm, nil
}

func (c *ClientImpl) Ping() error {
	resp, err := c.apiClient.GetVmmPing(c.context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not ping vm manager: %s", string(body))
	}

	ping := api.VmmPingResponse{}
	err = json.NewDecoder(resp.Body).Decode(&ping)
	if err != nil {
		return err
	}

	return nil
}

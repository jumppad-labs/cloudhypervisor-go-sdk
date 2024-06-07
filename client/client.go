package client

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/instruqt/cloudhypervisor-go-sdk/api"
)

type Option func(*ClientImpl) error

const (
	defaultSocket = "/tmp/cloud-hypervisor.sock"
	defaultURL    = "http://localhost/api/v1/"
)

type Client interface {
	Create(api.VmConfig) (*api.VmInfo, error)
	Boot() (*api.VmInfo, error)
	Pause() (*api.VmInfo, error)
	Resume() (*api.VmInfo, error)
	Snapshot(api.VmSnapshotConfig) error
	Restore(api.RestoreConfig) error
	Reboot() (*api.VmInfo, error)
	PowerButton() (*api.VmInfo, error)
	Shutdown() (*api.VmInfo, error)
	Delete() error
	Info() (*api.VmInfo, error)
	Ping() (*api.VmmPingResponse, error)
}

type ClientImpl struct {
	context    context.Context
	url        string
	socket     string
	unixClient *http.Client
	apiClient  *api.Client
}

func NewClient() (Client, error) {
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

func (c *ClientImpl) Create(api.VmConfig) (*api.VmInfo, error) {
	config := api.VmConfig{}
	_, err := c.apiClient.CreateVM(c.context, config)
	if err != nil {
		return nil, err
	}

	// TODO: check for 204

	info, err := c.Info()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Boot() (*api.VmInfo, error) {
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

func (c *ClientImpl) Pause() (*api.VmInfo, error) {
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

func (c *ClientImpl) Resume() (*api.VmInfo, error) {
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

func (c *ClientImpl) Snapshot(config api.VmSnapshotConfig) error {
	_, err := c.apiClient.PutVmSnapshot(c.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) Restore(config api.RestoreConfig) error {
	_, err := c.apiClient.PutVmRestore(c.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientImpl) Reboot() (*api.VmInfo, error) {
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

func (c *ClientImpl) PowerButton() (*api.VmInfo, error) {
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

func (c *ClientImpl) Shutdown() (*api.VmInfo, error) {
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
	_, err := c.apiClient.DeleteVM(c.context)
	if err != nil {
		return err
	}

	// TODO: check for 204

	return nil
}

func (c *ClientImpl) Info() (*api.VmInfo, error) {
	resp, err := c.apiClient.GetVmInfo(c.context)
	if err != nil {
		return nil, err
	}

	var info *api.VmInfo
	err = json.NewDecoder(resp.Body).Decode(info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (c *ClientImpl) Ping() (*api.VmmPingResponse, error) {
	resp, err := c.apiClient.GetVmmPing(c.context)
	if err != nil {
		return nil, err
	}

	var ping *api.VmmPingResponse
	err = json.NewDecoder(resp.Body).Decode(ping)
	if err != nil {
		return nil, err
	}

	return ping, nil
}

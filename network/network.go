package network

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/containernetworking/cni/libcni"
	"github.com/containernetworking/cni/pkg/invoke"
	"github.com/containernetworking/cni/pkg/version"
)

const (
	DefaultPluginPath string = "/opt/cni/bin"
	DefaultConfigPath string = "/etc/cni/net.d"
	DefaultNetnsPath  string = "/var/run/netns"
)

type Config struct {
	netnsPath  string
	pluginPath string
	configPath string
}

type NetworkManager struct {
	config Config
	cni    *libcni.CNIConfig
	sync.RWMutex
}

type Option func(*NetworkManager)

func New(opts ...Option) *NetworkManager {
	config := Config{
		netnsPath:  DefaultNetnsPath,
		pluginPath: DefaultPluginPath,
		configPath: DefaultConfigPath,
	}

	cni := libcni.NewCNIConfig(
		[]string{config.pluginPath},
		&invoke.DefaultExec{
			RawExec:       &invoke.RawExec{Stderr: os.Stderr},
			PluginDecoder: version.PluginDecoder{},
		},
	)

	return &NetworkManager{
		config: config,
		cni:    cni,
	}
}

func (n *NetworkManager) Create(ctx context.Context, id string) error {
	n.Lock()
	defer n.Unlock()

	cl, err := libcni.LoadNetworkConf(n.config.configPath, id)
	if err != nil {
		return err
	}

	netns := filepath.Join(n.config.netnsPath, id)
	ifName := "eth0"

	var cniArgs [][2]string
	var capabilityArgs map[string]any

	rt := &libcni.RuntimeConf{
		ContainerID:    id,
		NetNS:          netns,
		IfName:         ifName,
		Args:           cniArgs,
		CapabilityArgs: capabilityArgs,
	}

	result, err := n.cni.AddNetworkList(ctx, cl, rt)
	if err != nil {
		return err
	}

	// swallowing the error for now...
	_ = result.Print()

	return nil
}

func (n *NetworkManager) Destroy(ctx context.Context, id string) error {
	n.Lock()
	defer n.Unlock()

	cl, err := libcni.LoadNetworkConf(n.config.configPath, id)
	if err != nil {
		return err
	}

	netns := filepath.Join(n.config.netnsPath, id)
	ifName := "eth0"

	var cniArgs [][2]string
	var capabilityArgs map[string]any

	rt := &libcni.RuntimeConf{
		ContainerID:    id,
		NetNS:          netns,
		IfName:         ifName,
		Args:           cniArgs,
		CapabilityArgs: capabilityArgs,
	}

	err = n.cni.DelNetworkList(ctx, cl, rt)
	if err != nil {
		return err
	}

	return nil
}

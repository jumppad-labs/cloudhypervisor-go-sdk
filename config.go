package sdk

import (
	"encoding/json"

	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"
)

type Config struct {
	Kernel  KernelConfig    `json:"kernel"`
	Disks   []DiskConfig    `json:"disks"`
	Network []NetworkConfig `json:"networks"`
	Devices []DeviceConfig  `json:"devices"`

	CPU    CPUConfig    `json:"cpu"`
	Memory MemoryConfig `json:"memory"`

	Console *ConsoleConfig `json:"console"`
	Debug   *DebugConfig   `json:"debug"`
	Serial  *ConsoleConfig `json:"serial"`
	Vsock   *VsockConfig   `json:"vsock"`
}

func (c *Config) Validate() error {
	return nil
}

/*
Converting the Config to VmConfig by marshalling and unmarshalling
because we are basically wrapping the client.VmConfig and extending it.

If we are changing our config further away from the client.VmConfig,
we need to change this function.
*/
func (c *Config) ToVmConfig() (client.VmConfig, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return client.VmConfig{}, err
	}

	var vmConfig client.VmConfig
	err = json.Unmarshal(b, &vmConfig)
	if err != nil {
		return client.VmConfig{}, err
	}

	return vmConfig, nil
}

type KernelConfig struct {
	Args   *string `json:"args,omitempty"`
	Initrd *string `json:"initrd,omitempty"`
	Path   *string `json:"path,omitempty"`
}

type DiskConfig struct {
	ID       *string `json:"id,omitempty"`
	Path     string  `json:"path"`
	Readonly *bool   `json:"readonly,omitempty"`
	Shared   *bool
}

type NetworkConfig struct {
	ID            *string `json:"id,omitempty"`
	IpAddress     *string `json:"ip,omitempty"`
	MacAddress    *string `json:"mac,omitempty"`
	Mask          *string `json:"mask,omitempty"`
	HostInterface *string `json:"tap,omitempty"`
	Gateway       *string
}

type DeviceConfig struct {
	ID   *string `json:"id,omitempty"`
	Path string  `json:"path"`
}

type CPUConfig struct {
	BootVcpus int `json:"boot_vcpus"`
	MaxVcpus  int `json:"max_vcpus"`
}

type MemoryConfig struct {
	Size int64 `json:"size"`
}

type ConsoleMode string

const (
	ConsoleModeOff    ConsoleMode = "Off"
	ConsoleModeTty    ConsoleMode = "Tty"
	ConsoleModePty    ConsoleMode = "Pty"
	ConsoleModeFile   ConsoleMode = "File"
	ConsoleModeSocket ConsoleMode = "Socket"
	ConsoleModeNull   ConsoleMode = "Null"
)

type ConsoleConfig struct {
	File   *string     `json:"file,omitempty"`
	Mode   ConsoleMode `json:"mode"`
	Socket *string     `json:"socket,omitempty"`
}

type DebugConfig struct {
	File *string     `json:"file,omitempty"`
	Mode ConsoleMode `json:"mode"`
}

type VsockConfig struct {
	ID     *string `json:"id,omitempty"`
	Socket string  `json:"socket"`
}

type VMState string

const (
	VMStateCreated  VMState = "Created"
	VMStatePaused   VMState = "Paused"
	VMStateRunning  VMState = "Running"
	VMStateShutdown VMState = "Shutdown"
)

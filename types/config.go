package types

import (
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/api"
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

func (c Config) Convert() api.VmConfig {
	vmConfig := api.VmConfig{
		Cpus:    c.CPU.Convert(),
		Memory:  c.Memory.Convert(),
		Payload: c.Kernel.Convert(),
		Disks:   &[]api.DiskConfig{},
		Net:     &[]api.NetConfig{},
		Devices: &[]api.DeviceConfig{},
	}

	for _, disk := range c.Disks {
		*vmConfig.Disks = append(*vmConfig.Disks, disk.Convert())
	}

	for _, network := range c.Network {
		*vmConfig.Net = append(*vmConfig.Net, network.Convert())
	}

	for _, device := range c.Devices {
		*vmConfig.Devices = append(*vmConfig.Devices, device.Convert())
	}

	if c.Console != nil {
		vmConfig.Console = c.Console.Convert()
	}

	if c.Debug != nil {
		vmConfig.DebugConsole = c.Debug.Convert()
	}

	if c.Serial != nil {
		vmConfig.Serial = c.Serial.Convert()
	}

	if c.Vsock != nil {
		vmConfig.Vsock = c.Vsock.Convert()
	}

	return vmConfig
}

type KernelConfig struct {
	Args   *string `json:"args,omitempty"`
	Initrd *string `json:"initrd,omitempty"`
	Path   *string `json:"path,omitempty"`
}

func (c KernelConfig) Convert() api.PayloadConfig {
	return api.PayloadConfig{
		Cmdline:   c.Args,
		Initramfs: c.Initrd,
		Kernel:    c.Path,
	}
}

type DiskConfig struct {
	Direct   *bool   `json:"direct,omitempty"`
	ID       *string `json:"id,omitempty"`
	Path     string  `json:"path"`
	Readonly *bool   `json:"readonly,omitempty"`
}

func (c DiskConfig) Convert() api.DiskConfig {
	return api.DiskConfig{
		Direct: c.Direct,
		Id:     c.ID,
	}
}

type NetworkConfig struct {
	ID   *string `json:"id,omitempty"`
	IP   *string `json:"ip,omitempty"`
	MAC  *string `json:"mac,omitempty"`
	Mask *string `json:"mask,omitempty"`
	Tap  *string `json:"tap,omitempty"`
}

func (c NetworkConfig) Convert() api.NetConfig {
	return api.NetConfig{
		Id:   c.ID,
		Ip:   c.IP,
		Mac:  c.MAC,
		Mask: c.Mask,
		Tap:  c.Tap,
	}
}

type DeviceConfig struct {
	ID   *string `json:"id,omitempty"`
	Path string  `json:"path"`
}

func (c DeviceConfig) Convert() api.DeviceConfig {
	return api.DeviceConfig{
		Id: c.ID,
	}
}

type CPUConfig struct {
	BootVcpus int `json:"boot_vcpus"`
	MaxVcpus  int `json:"max_vcpus"`
}

func (c CPUConfig) Convert() *api.CpusConfig {
	return &api.CpusConfig{
		BootVcpus: c.BootVcpus,
		MaxVcpus:  c.MaxVcpus,
	}
}

type MemoryConfig struct {
	Size int64 `json:"size"`
}

func (c MemoryConfig) Convert() *api.MemoryConfig {
	return &api.MemoryConfig{
		Size: c.Size,
	}
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

func (c ConsoleConfig) Convert() *api.ConsoleConfig {
	return &api.ConsoleConfig{
		Mode:   api.ConsoleConfigMode(c.Mode),
		File:   c.File,
		Socket: c.Socket,
	}
}

type DebugConfig struct {
	File *string     `json:"file,omitempty"`
	Mode ConsoleMode `json:"mode"`
}

func (c DebugConfig) Convert() *api.DebugConsoleConfig {
	return &api.DebugConsoleConfig{
		Mode: api.DebugConsoleConfigMode(c.Mode),
		File: c.File,
	}
}

type VsockConfig struct {
	ID     *string `json:"id,omitempty"`
	Socket string  `json:"socket"`
}

func (c VsockConfig) Convert() *api.VsockConfig {
	return &api.VsockConfig{
		Id:     c.ID,
		Socket: c.Socket,
	}
}

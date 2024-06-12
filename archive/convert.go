package main

import "github.com/jumppad-labs/cloudhypervisor-go-sdk/client"

func VmInfoToVM(info *client.VmInfo) VM {
	return VM{
		Config: VmConfigToConfig(info.Config),
		State:  VMState(info.State),
	}
}

func VmConfigToConfig(vmConfig client.VmConfig) Config {
	config := Config{
		Kernel: KernelConfig{
			Args:   vmConfig.Payload.Cmdline,
			Initrd: vmConfig.Payload.Initramfs,
			Path:   vmConfig.Payload.Kernel,
		},
		Disks:   []DiskConfig{},
		Network: []NetworkConfig{},
		Devices: []DeviceConfig{},
		CPU: CPUConfig{
			BootVcpus: vmConfig.Cpus.BootVcpus,
			MaxVcpus:  vmConfig.Cpus.MaxVcpus,
		},
		Memory: MemoryConfig{
			Size: vmConfig.Memory.Size,
		},
	}

	if vmConfig.Disks != nil {
		for _, disk := range *vmConfig.Disks {
			config.Disks = append(config.Disks, DiskConfig{
				Direct:   disk.Direct,
				ID:       disk.Id,
				Path:     disk.Path,
				Readonly: disk.Readonly,
			})
		}
	}

	if vmConfig.Net != nil {
		for _, network := range *vmConfig.Net {
			config.Network = append(config.Network, NetworkConfig{
				ID:   network.Id,
				IP:   network.Ip,
				MAC:  network.Mac,
				Mask: network.Mask,
				Tap:  network.Tap,
			})
		}
	}

	if vmConfig.Devices != nil {
		for _, device := range *vmConfig.Devices {
			config.Devices = append(config.Devices, DeviceConfig{
				ID:   device.Id,
				Path: device.Path,
			})
		}
	}

	if vmConfig.Console != nil {
		config.Console = &ConsoleConfig{
			Mode:   ConsoleMode(vmConfig.Console.Mode),
			File:   vmConfig.Console.File,
			Socket: vmConfig.Console.Socket,
		}
	}

	if vmConfig.DebugConsole != nil {
		config.Debug = &DebugConfig{
			Mode: ConsoleMode(vmConfig.DebugConsole.Mode),
			File: vmConfig.DebugConsole.File,
		}
	}

	if vmConfig.Serial != nil {
		config.Serial = &ConsoleConfig{
			Mode:   ConsoleMode(vmConfig.Serial.Mode),
			File:   vmConfig.Serial.File,
			Socket: vmConfig.Serial.Socket,
		}
	}

	return config
}

func (c Config) Convert() client.VmConfig {
	vmConfig := client.VmConfig{
		Cpus:    c.CPU.Convert(),
		Memory:  c.Memory.Convert(),
		Payload: c.Kernel.Convert(),
	}

	disks := []client.DiskConfig{}
	for _, disk := range c.Disks {
		disks = append(disks, disk.Convert())
	}
	if len(disks) > 0 {
		vmConfig.Disks = &disks
	}

	networks := []client.NetConfig{}
	for _, network := range c.Network {
		networks = append(networks, network.Convert())
	}
	if len(networks) > 0 {
		vmConfig.Net = &networks
	}

	devices := []client.DeviceConfig{}
	for _, device := range c.Devices {
		devices = append(devices, device.Convert())
	}
	if len(devices) > 0 {
		vmConfig.Devices = &devices
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

func (c DiskConfig) Convert() client.DiskConfig {
	return client.DiskConfig{
		Direct:   c.Direct,
		Id:       nil,
		Path:     c.Path,
		Readonly: c.Readonly,
	}
}

func (c NetworkConfig) Convert() client.NetConfig {
	return client.NetConfig{
		Id:      nil,
		Ip:      c.IP,
		Mac:     c.MAC,
		Mask:    c.Mask,
		Tap:     c.Tap,
		HostMac: nil,
	}
}

func (c DeviceConfig) Convert() client.DeviceConfig {
	return client.DeviceConfig{
		Id: c.ID,
	}
}

func (c CPUConfig) Convert() *client.CpusConfig {
	return &client.CpusConfig{
		BootVcpus: c.BootVcpus,
		MaxVcpus:  c.MaxVcpus,
	}
}

func (c MemoryConfig) Convert() *client.MemoryConfig {
	return &client.MemoryConfig{
		Size: c.Size,
	}
}

func (c ConsoleConfig) Convert() *client.ConsoleConfig {
	return &client.ConsoleConfig{
		Mode:   client.ConsoleConfigMode(c.Mode),
		File:   c.File,
		Socket: c.Socket,
	}
}

func (c DebugConfig) Convert() *client.DebugConsoleConfig {
	return &client.DebugConsoleConfig{
		Mode: client.DebugConsoleConfigMode(c.Mode),
		File: c.File,
	}
}

func (c VsockConfig) Convert() *client.VsockConfig {
	return &client.VsockConfig{
		Id:     c.ID,
		Socket: c.Socket,
	}
}

func (c KernelConfig) Convert() client.PayloadConfig {
	return client.PayloadConfig{
		Cmdline:   c.Args,
		Initramfs: c.Initrd,
		Kernel:    c.Path,
	}
}

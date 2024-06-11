package types

import "github.com/jumppad-labs/cloudhypervisor-go-sdk/api"

func VmInfoToVM(info *api.VmInfo) VM {
	return VM{
		Config: VmConfigToConfig(info.Config),
		State:  VMState(info.State),
	}
}

func VmConfigToConfig(vmConfig api.VmConfig) Config {
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

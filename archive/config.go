package main

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

type KernelConfig struct {
	Args   *string `json:"args,omitempty"`
	Initrd *string `json:"initrd,omitempty"`
	Path   *string `json:"path,omitempty"`
}

type DiskConfig struct {
	Direct   *bool   `json:"direct,omitempty"`
	ID       *string `json:"id,omitempty"`
	Path     string  `json:"path"`
	Readonly *bool   `json:"readonly,omitempty"`
}

type NetworkConfig struct {
	ID   *string `json:"id,omitempty"`
	IP   *string `json:"ip,omitempty"`
	MAC  *string `json:"mac,omitempty"`
	Mask *string `json:"mask,omitempty"`
	Tap  *string `json:"tap,omitempty"`
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

type VM struct {
	Config Config  `json:"config"`
	State  VMState `json:"state"`
}

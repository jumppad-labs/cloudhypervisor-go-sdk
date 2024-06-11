package types

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

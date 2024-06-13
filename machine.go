package sdk

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/jumppad-labs/cloudhypervisor-go-sdk/client"
)

// TODO: set up networking
// TODO: handle signals
// TODO: set up vm/vmm logging
// TODO: set up vm/vmm metrics
// TODO: create overlayfs disk
// TODO: virtio-fs /var/lib/docker/overlay2 to mount docker images into guest

/*
# set the kernel boot args to use overlay-init
"boot_args": "console=ttyS0 reboot=k panic=1 pci=off overlay_root=vdb init=/sbin/overlay-init"
  - overlay_root: the disk that is the overlay root
  - init: override the default init program to set up the overlay filesystem

# create read only filesystem
sudo mkdir -p $MOUNTDIR/overlay/root $MOUNTDIR/overlay/work $MOUNTDIR/mnt $MOUNTDIR/rom
sudo cp files/overlay-init $MOUNTDIR/sbin/overlay-init
sudo mksquashfs $MOUNTDIR $SQUASHFS -noappend

https://github.com/cloud-hypervisor/cloud-hypervisor/blob/main/docs/custom-image.md
*/
type Option func(*MachineImpl) error

const (
	defaultSocket = "/tmp/cloud-hypervisor.sock"
	defaultURL    = "http://localhost/api/v1/"
)

type Machine interface {
	PID() (int, error)
	Start(ctx context.Context) error
	Pause(ctx context.Context) error
	Resume(ctx context.Context) error
	Snapshot(ctx context.Context, destination string) error
	Restore(ctx context.Context, source string) error
	Reboot(ctx context.Context) error
	PowerButton(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Wait(ctx context.Context) error
	Info(ctx context.Context) (*client.VmInfo, error)
}

type MachineImpl struct {
	context       context.Context
	client        *client.Client
	cmd           *exec.Cmd
	config        Config
	machineConfig client.VmConfig
	startOnce     sync.Once
	exitCh        chan struct{}
	fatalErr      error
	logger        *log.Logger
}

func NewMachine(ctx context.Context, config Config) (Machine, error) {
	logger := log.New(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	path, err := exec.LookPath("cloud-hypervisor")
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, path, "--api-socket", defaultSocket)

	unixClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", defaultSocket)
			},
		},
	}

	client, err := client.NewClient(defaultURL, client.WithHTTPClient(unixClient))
	if err != nil {
		return nil, err
	}

	// TODO: validate config
	// err = config.Validate()

	// TODO: convert config to vm config

	return &MachineImpl{
		context: ctx,
		client:  client,
		cmd:     cmd,
		config:  config,
		exitCh:  make(chan struct{}),
		logger:  logger,
	}, nil
}

func (m *MachineImpl) PID() (int, error) {
	if m.cmd == nil || m.cmd.Process == nil {
		return 0, fmt.Errorf("machine is not running")
	}

	select {
	case <-m.exitCh:
		return 0, fmt.Errorf("machine process has exited")
	default:
	}
	return m.cmd.Process.Pid, nil
}

func (m *MachineImpl) Start(ctx context.Context) error {
	alreadyStarted := true
	m.startOnce.Do(func() {
		m.logger.Debug("marking machine as started")
		alreadyStarted = false
	})
	if alreadyStarted {
		return fmt.Errorf("machine already started")
	}

	// start vmm
	err := m.startVMM()
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		waitErr := m.cmd.Wait()
		if waitErr != nil {
			errCh <- waitErr
		}

		close(errCh)
	}()

	err = m.waitForSocket(10*time.Second, errCh)
	if err != nil {
		m.fatalErr = err
		close(m.exitCh)
	}

	m.logger.Debug("vmm is ready")

	err = m.createVM()
	if err != nil {
		m.fatalErr = err
		close(m.exitCh)
	}

	err = m.bootVM()
	if err != nil {
		m.fatalErr = err
		close(m.exitCh)
	}

	return nil
}

func (m *MachineImpl) startVMM() error {
	err := m.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (m *MachineImpl) waitForSocket(timeout time.Duration, exitCh chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	ticker := time.NewTicker(10 * time.Millisecond)

	defer func() {
		cancel()
		ticker.Stop()
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exitCh:
			return err
		case <-ticker.C:
			if _, err := os.Stat(defaultSocket); err != nil {
				continue
			}

			if err := m.ping(); err != nil {
				continue
			}

			return nil
		}
	}
}

func (m *MachineImpl) createVM() error {
	resp, err := m.client.CreateVM(m.context, m.machineConfig)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not create vm: %s", string(body))
	}

	return nil
}

func (m *MachineImpl) bootVM() error {
	resp, err := m.client.BootVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not boot vm: %s", string(body))
	}

	return nil
}

func (m *MachineImpl) Pause(ctx context.Context) error {
	_, err := m.client.PauseVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404, 405
	return nil
}

func (m *MachineImpl) Resume(ctx context.Context) error {
	_, err := m.client.ResumeVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404, 405
	return nil
}

func (m *MachineImpl) Snapshot(ctx context.Context, destination string) error {
	config := client.VmSnapshotConfig{}
	_, err := m.client.PutVmSnapshot(m.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (m *MachineImpl) Restore(ctx context.Context, source string) error {
	config := client.RestoreConfig{}
	_, err := m.client.PutVmRestore(m.context, config)
	if err != nil {
		return err
	}

	return nil
}

func (m *MachineImpl) Reboot(ctx context.Context) error {
	_, err := m.client.RebootVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404, 405
	return nil
}

func (m *MachineImpl) PowerButton(ctx context.Context) error {
	_, err := m.client.PowerButtonVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404, 405

	return nil
}

func (m *MachineImpl) Shutdown(ctx context.Context) error {
	resp, err := m.client.ShutdownVM(m.context)
	if err != nil {
		return err
	}

	// TODO: check for 204, 404, 405
	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not shutdown vm: %s", string(body))
	}

	resp, err = m.client.ShutdownVMM(m.context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not shutdown vmm: %s", string(body))
	}

	return nil
}

func (m *MachineImpl) Delete() error {
	resp, err := m.client.DeleteVM(m.context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not delete vm: %s", string(body))
	}

	return nil
}

func (m *MachineImpl) ping() error {
	resp, err := m.client.GetVmmPing(m.context)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("could not ping vmm: %s", string(body))
	}

	_, err = client.ParseGetVmmPingResponse(resp)
	if err != nil {
		return err
	}

	return nil
}

func (m *MachineImpl) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.exitCh:
		return m.fatalErr
	}
}

func (m *MachineImpl) Version(ctx context.Context) (string, error) {
	resp, err := m.client.GetVmmPing(ctx)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("could not get vmm info: %s", string(body))
	}

	info, err := client.ParseGetVmmPingResponse(resp)
	if err != nil {
		return "", err
	}

	return info.JSON200.Version, nil
}

func (m *MachineImpl) Info(ctx context.Context) (*client.VmInfo, error) {
	resp, err := m.client.GetVmInfo(ctx)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("could not get vm info: %s", string(body))
	}

	info, err := client.ParseGetVmInfoResponse(resp)
	if err != nil {
		return nil, err
	}

	return info.JSON200, nil
}

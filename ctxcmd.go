package ctxcmd

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/net/context"
)

// Command represents an external command being prepared or run.
type Command struct {
	cmd *exec.Cmd
	ctx context.Context
}

// NewCommand bundles Context with a command to be run and returns an instance
// of *ctxtcmd.Command.
func NewCommand(ctx context.Context, cmd *exec.Cmd) *Command {
	return &Command{cmd, ctx}
}

// TerminationStatus is an error that arises when Command is prematurely
// terminated by a call from it's underlying Context.
type TerminationStatus struct {
	// ContextStatus carries information about how the Context was canceled.
	ContextStatus error

	// KillStatus carries the status of killing the underlying process. It
	// may be nil if killing the process succeeded.
	KillStatus error
}

func (t *TerminationStatus) Error() string {
	return fmt.Sprintf(
		"Context terminated with %v, process killed with %v",
		t.ContextStatus,
		t.KillStatus,
	)
}

// RunWithSignal starts the underlying exec.Cmd and waits until either the
// command finishes or the underlying Context is canceled. In the latter case
// the chosen OS signal is sent to the subprocess.
func (c *Command) RunWithSignal(onCancel os.Signal) error {
	if err := c.cmd.Start(); err != nil {
		return err
	}
	commandFinished := make(chan error)
	go func() { commandFinished <- c.cmd.Wait() }()
	select {
	case <-c.ctx.Done():
		break
	case err := <-commandFinished:
		return err
	}
	return &TerminationStatus{c.ctx.Err(), c.cmd.Process.Signal(onCancel)}
}

// Run starts the underlying exec.Cmd and waits until either the command
// finishes or the underlying Context is canceled. In the latter case the
// subprocess is immediately killed.
func (c *Command) Run() error {
	return c.RunWithSignal(os.Kill)
}

package ctxcmd

import (
	"bytes"
	"io"
	"os/exec"
)

// Wrapper provides a wrapper for exec.Cmd which captures it's standard output
// and standard error as in-memory buffers.
type Wrapper struct {
	Command *exec.Cmd
	Stdout  bytes.Buffer
	Stderr  bytes.Buffer
}

// NewWrapper is a constructor for a Wrapper struct which additionally alows the
// user to pass in an input stream.
func NewWrapper(cmd *exec.Cmd, stdin io.Reader) *Wrapper {
	wrapper := &Wrapper{Command: cmd}
	if stdin != nil {
		cmd.Stdin = stdin
	}
	cmd.Stdout = &wrapper.Stdout
	cmd.Stderr = &wrapper.Stderr
	return wrapper
}

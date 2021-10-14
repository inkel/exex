// exex provides a custom Cmd type that wraps exec.Cmd in a way that
// it will always capture standard error stream if execution fails
// with an exec.ExitError.
//
// The standard library exec package contains a very useful API to
// execute commands, however, the exec.Cmd.Run and exec.Cmd.Output
// methods behave differently in the case of a failed execution. In
// particular, exec.Cmd.Run will NOT populate exec.ExitError.Stderr in
// the case of failure, whereas exec.Cmd.Output will do. While this is
// explicitly noted in the exec package documentation, it is a source
// of confusion even for experienced users.
//
// Another issue with the standard library package is that if we use
// exec.Cmd.Output to only capture the error we will be incurring in
// unnecessary allocations because it uses a bytes.Buffer for
// capturing STDOUT, and also STDERR will be truncated. The wrappers
// defined in this library do not have this peformance penalization
// nor truncate any output.
package exex

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
)

// Cmd wraps exec.Cmd and represents an external command.
//
// As in the case of exec.Cmd, a Cmd cannot be reused after executed
// for the first time.
type Cmd exec.Cmd

// Command returns the Cmd struct to execute the named program with
// the given arguments.
//
// Refer to the exec.Command documentation for additional information.
func Command(name string, args ...string) *Cmd {
	return (*Cmd)(exec.Command(name, args...))
}

// CommandContext is like Command but the Cmd is associated with a
// context.
//
// Refer to the exec.Command documentation for additional information.
func CommandContext(ctx context.Context, name string, args ...string) *Cmd {
	return (*Cmd)(exec.CommandContext(ctx, name, args...))
}

// Run starts the command and waits for it to end.
//
// If the command executes successfully (e.g. exits with a zero
// status), it returns nil.
//
// If the command fails to execute, the error will be of type
// *exec.ExitError and it's always guaranteed that its Stderr property
// will have the contexts of the standard error stream, unless
// *Cmd.Stderr is specified.
//
// Refer to exec.Cmd.Run documentation for additional information.
func (c *Cmd) Run() error {
	var stderr *bytes.Buffer

	if c.Stderr == nil {
		stderr = bytes.NewBuffer(make([]byte, 0, 1024))
		c.Stderr = stderr
	}

	err := (*exec.Cmd)(c).Run()

	var exErr *exec.ExitError

	if stderr != nil && errors.As(err, &exErr) {
		exErr.Stderr = stderr.Bytes()
		return exErr
	}

	return err
}

// RunCommand wraps an *exec.Cmd into a Cmd and returns the result of
// calling *Cmd.Run.
func RunCommand(cmd *exec.Cmd) error {
	return (*Cmd)(cmd).Run()
}

// Run creates a Cmd and returns the result of executing *Cmd.Run.
func Run(cmd string, args ...string) error {
	return Command(cmd, args...).Run()
}

// RunContext creates a Cmd with the given context and returns the
// result of executing *Cmd.Run.
func RunContext(ctx context.Context, cmd string, args ...string) error {
	return CommandContext(ctx, cmd, args...).Run()
}

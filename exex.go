// Package exex provides a custom Cmd type that wraps exec.Cmd in a
// way that it will always capture standard error stream if execution
// fails with an exec.ExitError.
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
//
// In order to avoid importing both this package and os/exec, this
// package provides aliases for the variables and top-level functions
// os/exec provides.
package exex

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
)

// Cmd wraps exec.Cmd and represents an external command.
//
// As in the case of exec.Cmd, a Cmd cannot be reused after executed
// for the first time.
//
// Refer to the exec.Cmd documentation for information on all the
// functions this type provides except for Run, which is overwritten
// by this struct.
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

// Start starts the specified command but does not wait for it to
// complete.
func (c *Cmd) Start() error {
	if c.Stderr == nil {
		c.Stderr = bytes.NewBuffer(make([]byte, 0, 1024))
	}
	return (*exec.Cmd)(c).Start()
}

// Wait waits for the command to exit and waits for any copying to
// stdin or copying from stdout or stderr to complete.
func (c *Cmd) Wait() error {
	err := (*exec.Cmd)(c).Wait()

	var exErr *exec.ExitError

	if stderr, ok := c.Stderr.(*bytes.Buffer); ok && errors.As(err, &exErr) {
		exErr.Stderr = stderr.Bytes()
		return exErr
	}

	return err
}

// Output runs the command and returns its standard output. Any
// returned error will usually be of type *ExitError. If c.Stderr was
// nil, Output populates ExitError.Stderr.
func (c *Cmd) Output() ([]byte, error) { return (*exec.Cmd)(c).Output() }

// CombinedOutput runs the command and returns its combined standard
// output and standard error.
func (c *Cmd) CombinedOutput() ([]byte, error) { return (*exec.Cmd)(c).CombinedOutput() }

// StderrPipe returns a pipe that will be connected to the command's
// standard error when the command starts.
func (c *Cmd) StderrPipe() (io.ReadCloser, error) { return (*exec.Cmd)(c).StderrPipe() }

// StdinPipe returns a pipe that will be connected to the command's
// standard input when the command starts.
func (c *Cmd) StdinPipe() (io.WriteCloser, error) { return (*exec.Cmd)(c).StdinPipe() }

// StdoutPipe returns a pipe that will be connected to the command's
// standard output when the command starts.
func (c *Cmd) StdoutPipe() (io.ReadCloser, error) { return (*exec.Cmd)(c).StdoutPipe() }

// String returns a human-readable description of c
func (c *Cmd) String() string { return (*exec.Cmd)(c).String() }

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

// CommandError returns the error with the stderr log appended,
// if the error is a command exit error, and the stderr log exists.
func CommandError(err error, errMsg string) error {
	var exErr *ExitError
	if err != nil {
		if !errors.As(err, &exErr) {
			return fmt.Errorf("error converting error to exex.ExitError")
		}
		if exErr.Stderr == nil {
			return fmt.Errorf("%s (%w)", errMsg, err)
		}
		return fmt.Errorf("%s (%w)\n%s", errMsg, err, exErr.Stderr)
	}
	return nil
}

// Error is a type alias for exec.Error
type Error = exec.Error

// ExitError is a type alias for exec.ExitError
type ExitError = exec.ExitError

// ErrNotFound is an alias for exec.ErrNotFound, the error resulting
// if a path search failed to find an executable file.
var ErrNotFound = exec.ErrNotFound

// LookPath is an alias for exec.LookPath, searches for an executable
// named file in the directories named by the PATH environment
// variable. Refer to that package for additional information.
var LookPath = exec.LookPath

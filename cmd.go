// exex provides syntactic sugar around exec.Cmd to easily run
// commands and capturing stderr.
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
	"os/exec"
)

// Run executes the given command and waits for it to complete.
//
// It behaves in the same way as exec.Cmd.Run, with the given
// differfence that if the returned error is of type *exec.ExitError
// it will have exec.ExitError.Stderr populated.
func RunCommand(cmd *exec.Cmd) error {
	var stderr *bytes.Buffer

	if cmd.Stderr == nil {
		stderr = bytes.NewBuffer(make([]byte, 0, 1024))
		cmd.Stderr = stderr
	}

	err := cmd.Run()

	if exErr, ok := err.(*exec.ExitError); ok && stderr != nil {
		exErr.Stderr = stderr.Bytes()
		return exErr
	}

	return err
}

// Run calls exec.Command with the given arguments and run it with RunCommand.
func Run(cmd string, args ...string) error {
	return RunCommand(exec.Command(cmd, args...))
}

// Run calls exec.CommandContext with the given arguments and run it with RunCommand.
func RunContext(ctx context.Context, cmd string, args ...string) error {
	return RunCommand(exec.CommandContext(ctx, cmd, args...))
}

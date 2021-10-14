package exex_test

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/inkel/exex"
)

func ExampleCommand() {
	cmd := exex.Command("true")
	err := cmd.Run()

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

func ExampleCommandContext() {
	cmd := exex.Command("true")
	err := cmd.Run()

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

func ExampleCmd_Run() {
	err := exex.Command("foo").Run()

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

func ExampleRun() {
	err := exex.Run("foo")

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

func ExampleRunContext() {
	ctx := context.Background()
	err := exex.RunContext(ctx, "foo")

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

func ExampleRunCommand() {
	cmd := exec.Command("foo")
	err := exex.RunCommand(cmd)

	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		fmt.Printf("Captured stderr: %q\n", exErr.Stderr)
	} else {
		fmt.Printf("Expecting an *exec.ExitError, got %T: %[1]v\n", err)
	}
}

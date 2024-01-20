package exex_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.arcalot.io/assert"
	"go.arcalot.io/exex"
	"go.arcalot.io/log/v2"
	"os"
	"os/exec"
	"path"
	"testing"
)

const stderrMessage = "Yup, I'm broken"

func TestMain(m *testing.M) {
	logger := log.NewLogger(log.LevelDebug, log.NewBufferWriter())
	if o := os.Getenv("TEST_MAIN"); o != "" {
		_, err := fmt.Fprint(os.Stderr, "error:")
		if err != nil {
			logger.Errorf("main failed to print to stderr %v", err)
			os.Exit(1)
		}
		for _, m := range os.Args[1:] {
			_, err2 := fmt.Fprint(os.Stderr, " ", m)
			if err2 != nil {
				logger.Errorf("main failed to print to stderr %v", err2)
				os.Exit(1)
			}
		}
		os.Exit(1)
	}

	bench := os.Getenv("BENCHMARK")
	os.Clearenv()
	err := os.Setenv("TEST_MAIN", "error")
	if err != nil {
		logger.Errorf("error setting TEST_MAIN in system environment %v", err)
		os.Exit(1)
	}
	err = os.Setenv("BENCHMARK", bench)
	if err != nil {
		logger.Errorf("error setting BENCHMARK in system environment %v", err)
	}
	os.Exit(m.Run())
}

func assertErr(t *testing.T, err error, msg string) {
	assert.Error(t, err)
	var exErr *exec.ExitError
	assert.Equals(t, errors.As(err, &exErr), true)
	assert.Contains(t, string(exErr.Stderr), msg)
}

func TestRun(t *testing.T) {
	t.Run("command", func(t *testing.T) {
		err := exex.Run(os.Args[0])
		assertErr(t, err, "error:")
	})

	t.Run("command+args", func(t *testing.T) {
		err := exex.Run(os.Args[0], "foo", "bar")
		assertErr(t, err, "error: foo bar")
	})
}

func TestRunContext(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expecting error")
			}
			if r != "nil Context" {
				t.Fatalf("expecting nil context error, got %q", r)
			}
		}()

		err := exex.RunContext(nil, os.Args[0])
		assert.NoError(t, err)
	})

	t.Run("background", func(t *testing.T) {
		err := exex.RunContext(context.Background(), os.Args[0], "context")
		assertErr(t, err, "error: context")
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := exex.RunContext(ctx, os.Args[0], "context cancelled")
		if err == nil {
			t.Fatal("expecting error")
		}
		if ctx.Err() != err {
			t.Fatalf("expecting %v, got %v", ctx.Err(), err)
		}
	})
}

func TestCmd_RunCapture(t *testing.T) {
	fmt.Printf("%v\n", os.Args[0])
	cmd := exec.Command(os.Args[0], "capture", "stderr")
	err := exex.RunCommand(cmd)
	assertErr(t, err, "error: capture stderr")
}

func TestRunCommand(t *testing.T) {
	t.Run("capture", func(t *testing.T) {
		cmd := exec.Command(os.Args[0], "capture", "stderr")
		err := exex.RunCommand(cmd)
		assertErr(t, err, "error: capture stderr")
	})

	t.Run("custom stderr", func(t *testing.T) {
		var stderr bytes.Buffer
		cmd := exec.Command(os.Args[0], "capture", "stderr")
		cmd.Stderr = &stderr
		err := exex.RunCommand(cmd)
		if err == nil {
			t.Fatal("expecting error")
		}

		exErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("expecting *exec.ExitError, got %T", err)
		}
		if exErr.Stderr != nil {
			t.Errorf("expecting not captured stderr, got %q", exErr.Stderr)
		}

		exp := "error: capture stderr"
		if got := stderr.String(); got != exp {
			t.Errorf("expecting %q, got %q", exp, got)
		}
	})
}

var Stderr []byte

func benchmarkCaptureStderrStdlib(b *testing.B) {
	var exErr *exec.ExitError

	for i := 0; i < b.N; i++ {
		cmd := exec.Command(os.Args[0])
		cmd.Env = []string{"TEST_MAIN=error"}
		_, err := cmd.Output()
		exErr = err.(*exec.ExitError)
	}

	Stderr = exErr.Stderr
}

func benchmarkCaptureStderrExex(b *testing.B) {
	var exErr *exec.ExitError

	for i := 0; i < b.N; i++ {
		cmd := exec.Command(os.Args[0])
		cmd.Env = []string{"TEST_MAIN=error"}
		err := exex.RunCommand(cmd)
		exErr = err.(*exec.ExitError)
	}

	Stderr = exErr.Stderr
}

func BenchmarkCaptureStderr(b *testing.B) {
	switch os.Getenv("BENCHMARK") {
	case "stdlib":
		benchmarkCaptureStderrStdlib(b)
	case "exex":
		benchmarkCaptureStderrExex(b)
	default:
		b.Run("stdlib", benchmarkCaptureStderrStdlib)
		b.Run("exex", benchmarkCaptureStderrExex)
	}
}

func BenchmarkRun(b *testing.B) {
	var exErr *exec.ExitError

	for i := 0; i < b.N; i++ {
		exErr = exex.Run(os.Args[0]).(*exec.ExitError)
	}

	Stderr = exErr.Stderr
}

func BenchmarkRunContext(b *testing.B) {
	ctx := context.Background()

	var exErr *exec.ExitError

	for i := 0; i < b.N; i++ {
		exErr = exex.RunContext(ctx, os.Args[0]).(*exec.ExitError)
	}

	Stderr = exErr.Stderr
}

func BenchmarkRunCommand(b *testing.B) {
	var exErr *exec.ExitError

	for i := 0; i < b.N; i++ {
		exErr = exex.RunCommand(exec.Command(os.Args[0])).(*exec.ExitError)
	}

	Stderr = exErr.Stderr
}

func TestCmd_Run(t *testing.T) {
	t.Run("capture", func(t *testing.T) {
		err := exex.Command(os.Args[0], "capture", "stderr").Run()
		assertErr(t, err, "error: capture stderr")
	})

	t.Run("custom stderr", func(t *testing.T) {
		var stderr bytes.Buffer
		cmd := exex.Command(os.Args[0], "capture", "stderr")
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err == nil {
			t.Fatal("expecting error")
		}

		exErr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatalf("expecting *exec.ExitError, got %T", err)
		}
		if exErr.Stderr != nil {
			t.Errorf("expecting not captured stderr, got %q", exErr.Stderr)
		}

		exp := "error: capture stderr"
		if got := stderr.String(); got != exp {
			t.Errorf("expecting %q, got %q", exp, got)
		}
	})
}

func TestLookPathNotFound(t *testing.T) {
	nonExistentPath := "foobarbazquux"
	foundPath, err := exex.LookPath(nonExistentPath)
	assert.Error(t, err)
	assert.Equals(t, foundPath, "")
	_, ok := err.(*exex.Error)
	assert.Equals(t, ok, true)
}

func TestLookPathFound(t *testing.T) {
	bin := os.Args[0]
	t.Setenv("PATH", path.Dir(bin))
	binpath, err := exex.LookPath(path.Base(bin))
	assert.NoError(t, err)
	assert.Equals(t, binpath, bin)
}

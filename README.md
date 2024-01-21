# exec - Capture stderr in exec.Cmd.Run

`exex` provides a custom `Cmd` type that wraps [`exec.Cmd`](https://pkg.go.dev/os/exec#Cmd) in a way that it will always capture standard error stream if execution fails with an `exec.ExitError`.

The standard library `exec` package contains a very useful API to execute commands, however, the [exec.Cmd.Run](https://pkg.go.dev/os/exec#Cmd.Run) and [exec.Cmd.Output](https://pkg.go.dev/os/exec#Cmd.Output) methods behave differently in the case of a failed execution. In particular, `exec.Cmd.Run` will NOT populate `exec.ExitError.Stderr` in the case of failure, whereas `exec.Cmd.Output` will do. While this is explicitly noted in the exec package documentation, it is a source of confusion even for experienced users.

Another issue with the standard library package is that if we use `exec.Cmd.Output` to only capture the error we will be incurring in unnecessary allocations because it uses a `bytes.Buffer` for capturing STDOUT, and also STDERR will be truncated. The wrappers defined in this library do not have this peformance penalization nor truncate any output.

In order to avoid importing both this package and `os/exec`, this package provides aliases for the variables and top-level functions `os/exec` provides.

## License

See [LICENSE](./LICENSE), but basically, MIT.

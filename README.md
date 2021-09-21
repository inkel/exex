# exec - Capture stderr in exec.Cmd.Run

`exex` provides syntactic sugar around exec.Cmd to easily run commands and capturing STDERR.

The standard library `exec` package contains a very useful API to execute commands, however, the [exec.Cmd.Run](https://pkg.go.dev/os/exec#Cmd.Run) and [exec.Cmd.Output](https://pkg.go.dev/os/exec#Cmd.Output) methods behave differently in the case of a failed execution. In particular, `exec.Cmd.Run` will NOT populate `exec.ExitError.Stderr` in the case of failure, whereas `exec.Cmd.Output` will do. While this is explicitly noted in the exec package documentation, it is a source of confusion even for experienced users.

Another issue with the standard library package is that if we use `exec.Cmd.Output` to only capture the error we will be incurring in unnecessary allocations because it uses a `bytes.Buffer` for capturing STDOUT, and also STDERR will be truncated. The wrappers defined in this library do not have this peformance penalization nor truncate any output.


## Benchmarks

```shell
go test -benchmem -bench=. -benchtime=1000x
```

    goos: darwin
    goarch: amd64
    pkg: github.com/inkel/exex
    cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
    BenchmarkCaptureStderr/stdlib-12         	    1000	   3263849 ns/op	   35560 B/op	      49 allocs/op
    BenchmarkCaptureStderr/exex-12           	    1000	   3216620 ns/op	    2848 B/op	      39 allocs/op
    BenchmarkRun-12                          	    1000	   3532978 ns/op	    2904 B/op	      39 allocs/op
    BenchmarkRunContext-12                   	    1000	   3405547 ns/op	    3026 B/op	      41 allocs/op
    BenchmarkRunCommand-12                   	    1000	   3588594 ns/op	    2904 B/op	      39 allocs/op
    PASS
    ok  	github.com/inkel/exex	17.660s

As you can see, the number of bytes and allocations per operation is drastically improved in `exex`.

A better comparison of the difference in performance can be done using [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat):

```shell
# Generate benchmarks for stdlib
BENCHMARK=stdlib go test -benchmem -bench=CaptureStderr -benchtime=500x -count=10 > stdlib.txt

# Generate benchmarks for exex
BENCHMARK=exex go test -benchmem -bench=CaptureStderr -benchtime=500x -count=10 > exex.txt
```

```shell
# Compare the results
benchstat stdlib.txt exex.txt
```

    name              old time/op    new time/op    delta
    CaptureStderr-12    3.10ms ± 2%    3.11ms ± 2%     ~     (p=0.222 n=9+9)

    name              old alloc/op   new alloc/op   delta
    CaptureStderr-12    35.5kB ± 0%     2.9kB ± 2%  -91.91%  (p=0.000 n=9+10)

    name              old allocs/op  new allocs/op  delta
    CaptureStderr-12      49.0 ± 0%      39.0 ± 0%  -20.41%  (p=0.000 n=10+10)

If an image is worth a thousand words, how much worth is a benchmark?


## License

See [LICENSE](./LICENSE), but basically, MIT.

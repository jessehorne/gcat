gcat
===

An MIT-licensed [`cat`](https://www.gnu.org/software/coreutils/cat) implementation written in Go.

*Please Note*: This implementation isn't bug-free yet. There are bugs. I'm working on the ones I know about and I am probably unaware of others.. :)

# Overview

`gcat` is my attempt to further my understanding of Go and of the `cat` command. `gcat` attempts to be a complete implementation of `cat(1)` but I am noting here that much more testing is needed, especially concerning the presence of control characters and of `\r\n` combinations.

# Installation

```shell
go install github.com/jessehorne/gcat
gcat --help
```

# Benchmark

more coming soon...

```shell
dock@dock:~/source/goutils/gcat$ make benchmark
go test -bench=.
goos: linux
goarch: amd64
pkg: github.com/jessehorne/gcat
cpu: AMD Ryzen 5 2400G with Radeon Vega Graphics    
Benchmark_getOptions-8           1000000              1420 ns/op
Benchmark_parseArgs-8             174837              6443 ns/op
Benchmark_gcat-8                   35091             37814 ns/op
PASS
ok      github.com/jessehorne/gcat      5.364s

```

# Development

See `./Makefile`.

If you find a bug, please create an issue. <3

If you'd like to contribute, fork->issue->branch->pull request.

# License

See `./LICENSE`.
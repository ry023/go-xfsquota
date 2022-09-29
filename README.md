# go-xfsquota

A golang wrapper library for xfs_quota commandline tool.

## Overview

**Please note that this library will NOT work with Go binaries alone.**  

This library executes xfs_quota binary via the [os/exec](https://pkg.go.dev/os/exec) package.
Therefore, the xfs_quota binary must be deployed in the environment where this library is used.

Note also that a child process is created when exec binary.

## Installation

Installation can be done with a normal `go get`:

```
go get github.com/ry023/go-xfsquota
```

And, import package on your Go code.

```go
import xfsquota "github.com/ry023/go-xfsquota"
```

## Usage

Please see also [examples/](./examples)

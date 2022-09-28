# go-xfsquota-wrapper

A golang wrapper library for xfs_quota commandline tool.

**Please note that this library will NOT work with Go binaries alone.**  
This library executes xfs_quota binary via the [os/exec](https://pkg.go.dev/os/exec) package.
Therefore, the xfs_quota binary must be deployed in the environment where this library is used.

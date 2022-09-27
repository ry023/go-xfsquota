package main

import (
	"log"

	xq "github.com/ry023/go-xfsquota-wrapper"
)

const binaryPath = "/usr/bin/xfs_quota"
const mountPath = "/xfs_mount_root"

func main() {
	cli, err := xq.NewClient(binaryPath)
	if err != nil {
		log.Fatalf("caused error on creating client: %v", err)
	}

	gopt := xq.GlobalOption{
		Path: mountPath,
	}

	err = cli.Command(&gopt).SetupProjectWithId(100, "/path/to/dirtree_root", 0)
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}
}

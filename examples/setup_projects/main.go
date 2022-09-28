package main

import (
	"log"

	xq "github.com/ry023/go-xfsquota-wrapper"
)

const binaryPath = "/usr/bin/xfs_quota"
const fsPath = "/xfs_root"

func main() {
	cli, err := xq.NewClient(binaryPath)
	if err != nil {
		log.Fatalf("caused error on creating client: %v", err)
	}

  err = cli.Command(fsPath, nil).SetupProjectWithId(100, "/xfs_root/dirtree", 0)
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}
}

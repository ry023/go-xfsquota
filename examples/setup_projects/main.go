package main

import (
	"log"
  "fmt"

	xq "github.com/ry023/go-xfsquota-wrapper"
)

const binaryPath = "/usr/sbin/xfs_quota"
const fsPath = "/xfs_root"

func main() {
	cli, err := xq.NewClient(binaryPath)
	if err != nil {
		log.Fatalf("caused error on creating client: %v", err)
	}

	// Setup Project
	err = cli.Command(fsPath, nil).SetupProjectWithId(100, xq.ProjectCommandOption{Path: "/xfs_root/dir"})
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}

	// Limit pquota
	limitopt := xq.LimitCommandOption{
		Bsoft: 16384,
		Bhard: 16384,
		Isoft: 8192,
		Ihard: 8192,
	}
	err = cli.Command(fsPath, nil).LimitWithId(100, xq.QuotaTypeProject, xq.QuotaTargetTypeBlocks, limitopt)
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}

  // Report ProjectQuota
  rep, err := cli.Command(fsPath, nil).Report(xq.QuotaTypeProject, xq.QuotaTargetTypeBlocks, xq.ReportCommandOption{})
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}

  fmt.Printf("%+v\n", rep.ReportSets)
}

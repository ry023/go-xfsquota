package main

import (
	"context"
	"log"

	pp "github.com/k0kubun/pp/v3"
	xq "github.com/ry023/go-xfsquota-wrapper"
)

const binaryPath = "/usr/sbin/xfs_quota"
const fsPath = "/xfs_root"

func main() {
	cli, err := xq.NewClient(binaryPath)
	if err != nil {
		log.Fatalf("caused error on creating client: %v", err)
	}
	ctx := context.Background()

	// Setup Project
	err = cli.Command(fsPath, nil).SetupProjectWithId(ctx, 100, xq.ProjectCommandOption{Path: "/xfs_root/dir"})
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
	err = cli.Command(fsPath, nil).LimitWithId(ctx, 100, xq.QuotaTypeProject, xq.QuotaTargetTypeBlocks, limitopt)
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}

	// Report ProjectQuota
	rep, err := cli.Command(fsPath, nil).Report(ctx, xq.QuotaTypeProject, xq.QuotaTargetTypeBlocks, xq.ReportCommandOption{})
	if err != nil {
		log.Fatalf("caused error on execution: %v", err)
	}

	pp.Println(rep) // This is expected to output the following;
	// ```
	// &xfsquota.ReportResult{
	//   ReportSets: []xfsquota.ReportSet{
	//     xfsquota.ReportSet{
	//       QuotaType:       "Project",
	//       QuotaTargetType: "Blocks",
	//       MountPath:       "/xfs_root",
	//       DevicePath:      "",
	//       ReportValues:    []xfsquota.ReportValue{
	//         xfsquota.ReportValue{
	//           Id:    0,
	//           Used:  0,
	//           Soft:  0,
	//           Hard:  0,
	//           Grace: 0,
	//         },
	//         xfsquota.ReportValue{
	//           Id:    100,
	//           Used:  0,
	//           Soft:  16,
	//           Hard:  16,
	//           Grace: 0,
	//         },
	//       },
	//     },
	//   },
	// }
	// ```
}

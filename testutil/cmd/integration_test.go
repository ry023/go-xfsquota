package cmd_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ry023/go-xfsquota"
)

func TestNew(t *testing.T) {
	p, err := filepath.Abs("../../fake_xfs_quota")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := xfsquota.New(p); err != nil {
		t.Fatal(err)
	}
}

func TestProject(t *testing.T) {
	type args struct {
		ctx context.Context
		op  xfsquota.ProjectOpsType
		id  uint32
		opt xfsquota.ProjectCommandOption
	}

	type newCommandArgs struct {
		filesystemPath string
		globalOpt      *xfsquota.GlobalOption
	}

	tests := []struct {
		name           string
		newCommandArgs newCommandArgs
		args           args
		expectReport   []xfsquota.ReportSet
	}{
		{
			name: "setup project",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				op: xfsquota.ProjectSetupOps,
				id: 100,
				opt: xfsquota.ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectReport: []xfsquota.ReportSet{
				{
					QuotaType:       xfsquota.QuotaTypeProject,
					QuotaTargetType: xfsquota.QuotaTargetTypeBlocks,
					MountPath:       "/xfs_root",
					DevicePath:      "",
					ReportValues: []xfsquota.ReportValue{
						{},
						{Id: 100, Used: 0, Soft: 0, Hard: 0, Grace: 0},
					},
				},
			},
		},
		{
			name: "clear project",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				op: xfsquota.ProjectClearOps,
				id: 100,
				opt: xfsquota.ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectReport: []xfsquota.ReportSet{
				{
					QuotaType:       xfsquota.QuotaTypeProject,
					QuotaTargetType: xfsquota.QuotaTargetTypeBlocks,
					MountPath:       "/xfs_root",
					DevicePath:      "",
					ReportValues: []xfsquota.ReportValue{
						{},
					},
				},
			},
		},
	}
	p, err := filepath.Abs("../../fake_xfs_quota")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			q, err := xfsquota.New(p)
			if err != nil {
				t.Fatal(err)
			}
			c := q.Command(tt.newCommandArgs.filesystemPath, &xfsquota.GlobalOption{})
			if err := c.OperateProjectWithId(ctx, tt.args.op, tt.args.id, tt.args.opt); err != nil {
				t.Fatal(err)
			}
			r, err := c.Report(ctx, xfsquota.QuotaTypeProject, xfsquota.QuotaTargetTypeBlocks, xfsquota.ReportCommandOption{})
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tt.expectReport, r.ReportSets); diff != "" {
				t.Errorf("unexpected report (-want +got):\n%s", diff)
			}
		})
	}
}

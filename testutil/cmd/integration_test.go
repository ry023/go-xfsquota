package cmd_test

import (
	"context"
	"os/exec"
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
		op  xfsquota.ProjectOpsType
		id  uint32
		opt xfsquota.ProjectCommandOption
	}

	type newCommandArgs struct {
		filesystemPath string
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
		{
			name: "setup project2",
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
			name: "setup project3",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				op: xfsquota.ProjectSetupOps,
				id: 101,
				opt: xfsquota.ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir2",
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
						{Id: 101, Used: 0, Soft: 0, Hard: 0, Grace: 0},
					},
				},
			},
		},
	}
	p, err := filepath.Abs("../../fake_xfs_quota")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := exec.Command(p, "--fake-init").CombinedOutput(); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			q, err := xfsquota.New(p)
			if err != nil {
				t.Fatal(err)
			}
			c := q.Command(tt.newCommandArgs.filesystemPath, nil)
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

func TestLimit(t *testing.T) {
	type args struct {
		id        uint32
		quotaType xfsquota.QuotaType
		opt       xfsquota.LimitCommandOption
	}

	type newCommandArgs struct {
		filesystemPath string
	}

	tests := []struct {
		name           string
		newCommandArgs newCommandArgs
		args           args
		expectReport   []xfsquota.ReportSet
	}{
		{
			name: "limit projects",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				id:        100,
				quotaType: xfsquota.QuotaTypeProject,
				opt:       xfsquota.LimitCommandOption{Bsoft: 1024, Bhard: 1024, Isoft: 2048, Ihard: 2048, Rtbsoft: 512, Rtbhard: 512},
			},
			expectReport: []xfsquota.ReportSet{
				{
					QuotaType:       xfsquota.QuotaTypeProject,
					QuotaTargetType: xfsquota.QuotaTargetTypeBlocks,
					MountPath:       "/xfs_root",
					DevicePath:      "",
					ReportValues: []xfsquota.ReportValue{
						{},
						{Id: 100, Used: 0, Soft: 1024, Hard: 1024, Grace: 0},
					},
				},
			},
		},
		{
			name: "limit projects2",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				id:        101,
				quotaType: xfsquota.QuotaTypeProject,
				opt:       xfsquota.LimitCommandOption{Bsoft: 2048, Bhard: 2048, Isoft: 0, Ihard: 0, Rtbsoft: 0, Rtbhard: 0},
			},
			expectReport: []xfsquota.ReportSet{
				{
					QuotaType:       xfsquota.QuotaTypeProject,
					QuotaTargetType: xfsquota.QuotaTargetTypeBlocks,
					MountPath:       "/xfs_root",
					DevicePath:      "",
					ReportValues: []xfsquota.ReportValue{
						{},
						{Id: 100, Used: 0, Soft: 1024, Hard: 1024, Grace: 0},
						{Id: 101, Used: 0, Soft: 2048, Hard: 2048, Grace: 0},
					},
				},
			},
		},
	}
	p, err := filepath.Abs("../../fake_xfs_quota")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := exec.Command(p, "--fake-init").CombinedOutput(); err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			q, err := xfsquota.New(p)
			if err != nil {
				t.Fatal(err)
			}
			c := q.Command(tt.newCommandArgs.filesystemPath, nil)
			if err := c.LimitWithId(ctx, tt.args.id, tt.args.quotaType, tt.args.opt); err != nil {
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

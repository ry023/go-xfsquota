package xfsquota

import (
	"context"
	"testing"
)

func TestCommand_OperateProjectWithId(t *testing.T) {
	type args struct {
		ctx context.Context
		op  ProjectOpsType
		id  uint32
		opt ProjectCommandOption
	}

	type newCommandArgs struct {
		filesystemPath string
		globalOpt      *GlobalOption
	}

	tests := []struct {
		name           string
		newCommandArgs newCommandArgs
		args           args
		expectErr      bool
		expectBinArgs  []string
	}{
		{
			name: "setup project",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				op: ProjectSetupOps,
				id: 100,
				opt: ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "project -d 5 -p /xfs_root/dir -s 100", "/xfs_root"},
		},
		{
			name: "clear project",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				op: ProjectClearOps,
				id: 100,
				opt: ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "project -d 5 -p /xfs_root/dir -C 100", "/xfs_root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background() // Default
			}

			bin := &MockBinary{ExpectedArgs: tt.expectBinArgs}
			cmd := NewCommand(bin, tt.newCommandArgs.filesystemPath, nil, tt.newCommandArgs.globalOpt)

			err := cmd.OperateProjectWithId(tt.args.ctx, tt.args.op, tt.args.id, tt.args.opt)

			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

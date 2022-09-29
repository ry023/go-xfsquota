package xfsquota

import (
	"context"
	"testing"
)

func TestCommand_LimitWithId(t *testing.T) {
	type args struct {
		ctx       context.Context
		id        uint32
		quotaType QuotaType
		opt       LimitCommandOption
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
			name: "limit project's block",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				id:        100,
				quotaType: QuotaTypeProject,
				opt:       LimitCommandOption{Bsoft: 1024, Bhard: 1024, Isoft: 2048, Ihard: 2048, Rtbsoft: 512, Rtbhard: 512},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "limit -p bsoft=1024 bhard=1024 isoft=2048 ihard=2048 rtbsoft=512 rtbhard=512 100", "/xfs_root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background() // Default
			}

			bin := &MockBinary{ExpectedArgs: tt.expectBinArgs}
			cmd := NewCommand(bin, tt.newCommandArgs.filesystemPath, tt.newCommandArgs.globalOpt)

			err := cmd.LimitWithId(tt.args.ctx, tt.args.id, tt.args.quotaType, tt.args.opt)

			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

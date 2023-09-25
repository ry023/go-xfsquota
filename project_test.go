package xfsquota

import (
	"context"
	"testing"
)

func TestCommand_SetupDirectoryTree(t *testing.T) {
	type args struct {
		ctx    context.Context
		projid uint32
		opt    ProjectCommandOption
	}

	tests := []struct {
		name          string
		args          args
		expectErr     bool
		expectBinArgs []string
	}{
		{
			args: args{
				projid: 100,
				opt: ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "project -d 5 -p /xfs_root/dir -s 100", "/xfs_root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background() // Default
			}

			bin := &mockBinary{ExpectedArgs: tt.expectBinArgs}
			cmd := NewCommand(bin, "/xfs_root", &GlobalOption{})

			err := cmd.SetupDirectoryTree(tt.args.ctx, tt.args.projid, tt.args.opt)
			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

func TestCommand_ClearDirectoryTree(t *testing.T) {
	type args struct {
		ctx    context.Context
		projid uint32
		opt    ProjectCommandOption
	}

	tests := []struct {
		name          string
		args          args
		expectErr     bool
		expectBinArgs []string
	}{
		{
			args: args{
				projid: 100,
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

			bin := &mockBinary{ExpectedArgs: tt.expectBinArgs}
			cmd := NewCommand(bin, "/xfs_root", &GlobalOption{})

			err := cmd.ClearDirectoryTree(tt.args.ctx, tt.args.projid, tt.args.opt)
			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

func TestCommand_CheckDirectoryTree(t *testing.T) {
	type args struct {
		ctx    context.Context
		projid uint32
		opt    ProjectCommandOption
	}

	tests := []struct {
		name          string
		args          args
		expectErr     bool
		expectBinArgs []string
	}{
		{
			args: args{
				projid: 100,
				opt: ProjectCommandOption{
					Depth: 5,
					Path:  "/xfs_root/dir",
				},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "project -d 5 -p /xfs_root/dir -c 100", "/xfs_root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background() // Default
			}

			bin := &mockBinary{ExpectedArgs: tt.expectBinArgs}
			cmd := NewCommand(bin, "/xfs_root", &GlobalOption{})

			err := cmd.CheckDirectoryTree(tt.args.ctx, tt.args.projid, tt.args.opt)
			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

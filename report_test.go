package xfsquota

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommand_ReportWithId(t *testing.T) {
	type args struct {
		ctx             context.Context
		quotaType       QuotaType
		quotaTargetType QuotaTargetType
		id              uint32
		opt             ReportCommandOption
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
			name: "report project's blocks",
			newCommandArgs: newCommandArgs{
				filesystemPath: "/xfs_root",
			},
			args: args{
				id:              100,
				quotaType:       QuotaTypeProject,
				quotaTargetType: QuotaTargetTypeBlocks,
				opt:             ReportCommandOption{},
			},
			expectErr:     false,
			expectBinArgs: []string{"-x", "-c", "report -p -b -N", "/xfs_root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.ctx == nil {
				tt.args.ctx = context.Background() // Default
			}

			out := []byte(`#0                   0          0          0     00 [--------]
#100                12       1024      20480     05 [--------]
#200                 0          4      10240     00 [--------]`)

			bin := &MockBinary{ExpectedArgs: tt.expectBinArgs, Out: out}
			cmd := NewCommand(bin, tt.newCommandArgs.filesystemPath, tt.newCommandArgs.globalOpt)

			_, err := cmd.Report(tt.args.ctx, tt.args.quotaType, tt.args.quotaTargetType, tt.args.opt)

			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			bin.AssertArgs(t)
		})
	}
}

func Test_parseHeadlessReportOutput(t *testing.T) {
	type args struct {
		stdout          []byte
		quotaType       QuotaType
		quotaTargetType QuotaTargetType
		mountPath       string
		devicePath      string
	}

	tests := []struct {
		name         string
		args         args
		expectResult *ReportResult
		expectErr    bool
	}{
		{
			name: "report -p",

			args: args{
				quotaType:       QuotaTypeProject,
				quotaTargetType: QuotaTargetTypeBlocks,
				mountPath:       "/path/to/mount",
				devicePath:      "/dev/sdc",
				stdout: []byte(
					`#0                   0          0          0     00 [--------]
#100                12       1024      20480     05 [--------]
#200                 0          4      10240     00 [--------]`,
				),
			},

			expectResult: &ReportResult{
				ReportSets: []ReportSet{
					{
						QuotaType:       QuotaTypeProject,
						QuotaTargetType: QuotaTargetTypeBlocks,
						MountPath:       "/path/to/mount",
						DevicePath:      "/dev/sdc",
						ReportValues: []ReportValue{
							{Id: 0, Used: 0, Soft: 0, Hard: 0, Grace: 0},
							{Id: 100, Used: 12, Soft: 1024, Hard: 20480, Grace: 5},
							{Id: 200, Used: 0, Soft: 4, Hard: 10240, Grace: 0},
						},
					},
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHeadlessReportOutput(tt.args.stdout, tt.args.quotaType, tt.args.quotaTargetType, tt.args.mountPath, tt.args.devicePath)

			if (err != nil) != tt.expectErr {
				t.Errorf("Unexpected error :%v", err)
			}

			if diff := cmp.Diff(tt.expectResult, result); diff != "" {
				t.Errorf("Expected result and actual result do not match:\n%s", diff)
			}
		})
	}
}

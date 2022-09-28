package xfsquota

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

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

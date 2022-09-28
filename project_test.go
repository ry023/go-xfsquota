package xfsquota

import (
	"testing"
)

func TestProjectCommandOption_SubCommandString(t *testing.T) {
	tests := []struct {
		name           string
		option         ProjectCommandOption
		expectedOutput string
	}{
		{
			name: "full subcommand",
			option: ProjectCommandOption{
				Depth:     10,
				Path:      "/path/to/projects",
				Operation: ProjectSetupOpts,
				Id:        []uint32{1, 2, 3, 100},
				Name:      []string{"projectA", "projectB", "projectC"},
			},
			expectedOutput: "project -d 10 -p /path/to/projects -s 1 2 3 100 projectA projectB projectC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := tt.option.SubCommandString()
			if output != tt.expectedOutput {
				t.Errorf("Expected output not match.\nexpected output:\n\t`%s`\nactual:\n\t`%s`", tt.expectedOutput, output)
			}
		})
	}
}

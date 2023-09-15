package xfsquota

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/cli/safeexec"
)

type BinaryExecuter interface {
	Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, args ...string) error
	Validate() error
}

type Binary struct {
	// The path to xfs_quota binary
	Path string
}

func (b *Binary) Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, args ...string) error {
	e, err := safeexec.LookPath(b.Path)
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, e, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (b *Binary) Validate() error {
	// Check file existence
	if _, err := os.Stat(b.Path); err != nil {
		return err
	}

	return nil
}

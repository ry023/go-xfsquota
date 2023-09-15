//go:build linux
package xfsquota

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type BinaryExecuter interface {
	Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, args ...string) error
	Validate() error
}

type XfsQuotaBinary struct {
	// The path to xfs_quota binary
	Path string
}

func (b *XfsQuotaBinary) Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, b.Path, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func (b *XfsQuotaBinary) Validate() error {
	// Check file existence
	if _, err := os.Stat(b.Path); err != nil {
		return err
	}

	return nil
}

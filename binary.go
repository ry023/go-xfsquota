package xfsquota

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type BinaryExecuter interface {
	Execute(args ...string) ([]byte, []byte, error)
	Validate() error
}

type XfsQuotaBinary struct {
	// The path to xfs_quota binary
	Path string
}

func (b *XfsQuotaBinary) Execute(args ...string) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(b.Path, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	stdoutBytes, err := io.ReadAll(&stdout)
	if err != nil {
		return nil, nil, err
	}

	stderrBytes, err := io.ReadAll(&stderr)
	if err != nil {
		return nil, nil, err
	}

	return stdoutBytes, stderrBytes, nil
}

func (b *XfsQuotaBinary) Validate() error {
	// Check file existence
	if _, err := os.Stat(b.Path); err != nil {
		return err
	}

	return nil
}

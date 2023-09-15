package xfsquota

import (
	"context"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type MockBinary struct {
	ExpectedArgs []string
	ActualArgs   []string
	Out          []byte
}

func (m *MockBinary) Execute(ctx context.Context, stdout io.Writer, stderr io.Writer, args ...string) error {
	m.ActualArgs = args

	if len(m.Out) > 0 {
		stdout.Write(m.Out)
	}

	return nil
}

func (m *MockBinary) Validate() error {
	return nil
}

func (m *MockBinary) AssertArgs(t *testing.T) {
	if diff := cmp.Diff(m.ExpectedArgs, m.ActualArgs); diff != "" {
		t.Errorf("Args on mary exception not match!\n%v`", diff)
	}
}

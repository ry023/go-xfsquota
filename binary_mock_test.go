package xfsquota

import (
	"context"
	"io"
	"reflect"
	"testing"
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
	if !reflect.DeepEqual(m.ActualArgs, m.ExpectedArgs) {
		t.Errorf("Args on mary exection not match! expected is `%v` but actual `%v`", m.ExpectedArgs, m.ActualArgs)
	}
}

package cmd_test

import (
	"path/filepath"
	"testing"

	"github.com/ry023/go-xfsquota"
)

func TestNew(t *testing.T) {
	p, err := filepath.Abs("../../fake_xfs_quota")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := xfsquota.New(p); err != nil {
		t.Fatal(err)
	}
}

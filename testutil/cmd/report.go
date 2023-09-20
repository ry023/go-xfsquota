package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const dbPrefix = "fake_xfs_quota_db_"

func cmdReport(w io.Writer, args []string, mountPath string) error {
	p := storeJSONName(mountPath)
	if _, err := os.Stat(p); err != nil {
		return fmt.Errorf("fake_xfs_quota: %w", err)
	}
	b, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("fake_xfs_quota: %w", err)
	}
	var q Quota
	if err := json.Unmarshal(b, &q); err != nil {
		return fmt.Errorf("fake_xfs_quota: %w", err)
	}
	_, _ = fmt.Fprintf(w, "#%d%10d%10d%10d       00 [--------]\n", 0, 0, 0, 0) // default #0
	for _, p := range q.Projects {
		if len(p.Paths) == 0 && p.Bhard == 0 && p.Bsoft == 0 && p.Ihard == 0 && p.Isoft == 0 {
			continue
		}
		_, _ = fmt.Fprintf(w, "#%d%10d%10d%10d       00 [--------]\n", p.ID, 0, p.Bsoft, p.Bhard)
	}
	return nil
}

func storeJSONName(mountPath string) string {
	return fmt.Sprintf("%s%s.json", dbPrefix, normalizePath(mountPath))
}

func normalizePath(path string) string {
	return strings.ReplaceAll(path, string(filepath.Separator), "--")
}

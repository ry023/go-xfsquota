package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

func cmdLimit(w io.Writer, args []string, mountPath string) error {
	p := storeJSONName(mountPath)
	var q Quota
	if _, err := os.Stat(p); err != nil {
		q = Quota{MountPath: mountPath, Projects: map[uint32]Project{}}
	} else {
		b, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("fake_xfs_quota: %w", err)
		}
		if err := json.Unmarshal(b, &q); err != nil {
			return fmt.Errorf("fake_xfs_quota: %w", err)
		}
	}
	if !slices.Contains(args, "-p") {
		return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
	}
	var (
		id    uint32
		bsoft uint32
		bhard uint32
		isoft uint32
		ihard uint32
	)
	for _, a := range args {
		switch {
		case a == "-p":
			continue
		case strings.Contains(a, "="):
			kv := strings.Split(a, "=")
			if len(kv) != 2 {
				return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
			}
			v64, err := strconv.ParseUint(kv[1], 10, 32)
			if err != nil {
				return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
			}
			v := uint32(v64)
			switch kv[0] {
			case "bsoft":
				bsoft = v
			case "bhard":
				bhard = v
			case "isoft":
				isoft = v
			case "ihard":
				ihard = v
			}
		default:
			id64, err := strconv.ParseUint(a, 10, 32)
			if err != nil {
				return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
			}
			id = uint32(id64)
		}
	}
	if id == 0 {
		return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
	}
	proj, ok := q.Projects[id]
	if !ok {
		proj = Project{ID: id}
	}
	proj.Bsoft = bsoft
	proj.Bhard = bhard
	proj.Isoft = isoft
	proj.Ihard = ihard
	if len(proj.Paths) == 0 && proj.Bhard == 0 && proj.Bsoft == 0 && proj.Ihard == 0 && proj.Isoft == 0 {
		delete(q.Projects, id)
	} else {
		q.Projects[id] = proj
	}
	f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("fake_xfs_quota: %w", err)
	}
	defer f.Close()
	if err := json.NewEncoder(f).Encode(&q); err != nil {
		return fmt.Errorf("fake_xfs_quota: %w", err)
	}
	return nil
}

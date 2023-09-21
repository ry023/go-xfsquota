package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/ry023/go-xfsquota"
)

func cmdProject(w io.Writer, args []string, mountPath string) error {
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
	var (
		id    uint32
		path  string
		ptype xfsquota.ProjectOpsType
	)
	i := 0
	for {
		if i > len(args)-1 {
			break
		}
		a := args[i]
		switch {
		case slices.Contains([]string{"-p", "-g", "-u"}, a):
			i++
			continue
		case a == "-s":
			// setup
			ptype = xfsquota.ProjectSetupOps
			i++
			continue
		case a == "-C":
			// clear
			ptype = xfsquota.ProjectClearOps
			i++
			continue
		case a == "-c":
			// check
			ptype = xfsquota.ProjectCheckOps
			i++
			continue
		case strings.HasPrefix(a, "-"):
			i++
			i++
			continue
		}
		id64, err := strconv.ParseUint(a, 10, 32)
		if err != nil {
			if path != "" {
				return fmt.Errorf("failed to parse args: %v => id:%v, path:%v", args, id, path)
			}
			path = a
		} else {
			id = uint32(id64)
		}
		i++
	}
	if id == 0 || path == "" {
		return fmt.Errorf("failed to parse args: %v => id:%v, path:%v", args, id, path)
	}
	proj, ok := q.Projects[id]
	if !ok {
		proj = Project{ID: id}
	}
	switch ptype {
	case xfsquota.ProjectSetupOps:
		if !slices.Contains(proj.Paths, path) {
			proj.Paths = append(proj.Paths, path)
		}
	case xfsquota.ProjectClearOps:
		proj.Paths = slices.DeleteFunc(proj.Paths, func(p string) bool {
			return p == path
		})
	case xfsquota.ProjectCheckOps:
	}
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

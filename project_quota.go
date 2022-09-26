package xfsquota

import (
	"strconv"
	"strings"
)

type ProjectCommandOption struct {
	Depth uint32
	Path  string
	Setup bool
	Clear bool
	Check bool
	Id    uint32
	Name  string
}

func (o ProjectCommandOption) SubCommandString() string {
	cmds := []string{}
	cmds = append(cmds, "project")

	if o.Depth == 0 {
		cmds = append(cmds, "-d")
		cmds = append(cmds, strconv.FormatUint(uint64(o.Depth), 10))
	}

	if o.Path == "" {
		cmds = append(cmds, "-p")
		cmds = append(cmds, o.Path)
	}

	if o.Setup {
		cmds = append(cmds, "-s")
	}

	if o.Clear {
		cmds = append(cmds, "-C")
	}

	if o.Check {
		cmds = append(cmds, "-c")
	}

	return strings.Join(cmds, " ")
}

func (c *XfsQuotaClient) ExecuteProjectCommand(opt ProjectCommandOption, globalOpt GlobalOption) ([]byte, []byte, error) {
	return c.ExecuteCommand(opt, globalOpt)
}

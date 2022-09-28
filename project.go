package xfsquota

import (
	"strconv"
	"strings"
)

type projectCommandArgs struct {
	// Setup/Clear/Check operation
	operation ProjectOpsType
	// Project ID to target
	id []uint32
	// Project name to target
	name []string

	opt ProjectCommandOption
}

type ProjectCommandOption struct {
	// Equeal to "-d" flag on commandline.
	// This option allows to limit recursion level when processing project directories
	Depth uint32
	// Equeal to "-p" flag on commandline.
	// This option allows to specify project paths at command line ( instead of /etc/projects ).
	Path string
}

// Equeal to "-sCc" flag on commandline.
type ProjectOpsType string

const (
	ProjectSetupOps = ProjectOpsType("Setup")
	ProjectClearOps = ProjectOpsType("Clear")
	ProjectCheckOps = ProjectOpsType("Check")
)

// Build 'project' subcommand
//
// format:
//   project [ -cCs [ -d depth ] [ -p path ] id | name ]
func (o projectCommandArgs) subCommandString() string {
	cmds := []string{}
	cmds = append(cmds, "project")

	if o.opt.Depth != 0 {
		cmds = append(cmds, "-d")
		cmds = append(cmds, strconv.FormatUint(uint64(o.opt.Depth), 10))
	}

	if o.opt.Path != "" {
		cmds = append(cmds, "-p")
		cmds = append(cmds, o.opt.Path)
	}

	switch o.operation {
	case ProjectSetupOps:
		cmds = append(cmds, "-s")
	case ProjectClearOps:
		cmds = append(cmds, "-C")
	case ProjectCheckOps:
		cmds = append(cmds, "-c")
	}

	for _, id := range o.id {
		cmds = append(cmds, strconv.FormatUint(uint64(id), 10))
	}

	for _, name := range o.name {
		cmds = append(cmds, name)
	}

	return strings.Join(cmds, " ")
}

func (c *Command) OperateProjectWithId(op ProjectOpsType, id uint32, opt ProjectCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = &projectCommandArgs{
		id:        []uint32{id},
		operation: op,
		opt:       opt,
	}
	return c.Execute()
}

func (c *Command) SetupProjectWithId(id uint32, opt ProjectCommandOption) error {
	return c.OperateProjectWithId(ProjectSetupOps, id, opt)
}

func (c *Command) ClearProjectWithId(id uint32, opt ProjectCommandOption) error {
	return c.OperateProjectWithId(ProjectClearOps, id, opt)
}

package xfsquota

import (
	"strconv"
	"strings"
)

type ProjectCommandOption struct {
	// Equeal to "-d" flag on commandline.
	// This option allows to limit recursion level when processing project directories
	Depth uint32
	// Equeal to "-p" flag on commandline.
	// This option allows to specify project paths at command line ( instead of /etc/projects ).
	Path string
	// Setup/Clear/Check operation
	Operation ProjectOptsType
	// Project ID to target
	Id []uint32
	// Project name to target
	Name []string
}

// Equeal to "-sCc" flag on commandline.
type ProjectOptsType string

const (
	ProjectSetupOpts = ProjectOptsType("Setup")
	ProjectClearOpts = ProjectOptsType("Clear")
	ProjectCheckOpts = ProjectOptsType("Check")
)

// Build 'project' subcommand
//
// format:
//   project [ -cCs [ -d depth ] [ -p path ] id | name ]
func (o ProjectCommandOption) SubCommandString() string {
	cmds := []string{}
	cmds = append(cmds, "project")

	if o.Depth != 0 {
		cmds = append(cmds, "-d")
		cmds = append(cmds, strconv.FormatUint(uint64(o.Depth), 10))
	}

	if o.Path != "" {
		cmds = append(cmds, "-p")
		cmds = append(cmds, o.Path)
	}

	switch o.Operation {
	case ProjectSetupOpts:
		cmds = append(cmds, "-s")
	case ProjectClearOpts:
		cmds = append(cmds, "-C")
	case ProjectCheckOpts:
		cmds = append(cmds, "-c")
	}

	for _, id := range o.Id {
		cmds = append(cmds, strconv.FormatUint(uint64(id), 10))
	}

	for _, name := range o.Name {
		cmds = append(cmds, name)
	}

	return strings.Join(cmds, " ")
}

func (c *Command) Project(opt ProjectCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.SubOpt = opt
	return c.Execute()
}

func (c *Command) SetupProjectWithId(id uint32, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Operation: ProjectSetupOpts,
		Path:      path,
		Depth:     depth,
		Id:        []uint32{id},
	}

	return c.Project(opt)
}

func (c *Command) SetupProjectWithName(name string, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Operation: ProjectSetupOpts,
		Path:      path,
		Depth:     depth,
		Name:      []string{name},
	}

	return c.Project(opt)
}

func (c *Command) ClearProjectWithId(id uint32, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Operation: ProjectClearOpts,
		Path:      path,
		Depth:     depth,
		Id:        []uint32{id},
	}

	return c.Project(opt)
}

func (c *Command) ClearProjectWithName(name string, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Operation: ProjectClearOpts,
		Path:      path,
		Depth:     depth,
		Name:      []string{name},
	}

	return c.Project(opt)
}

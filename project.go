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
	// Equeal to "-s" flag on commandline.
	Setup bool
	// Equeal to "-C" flag on commandline.
	Clear bool
	// Equeal to "-c" flag on commandline.
	Check bool
	// Project ID to target
	Id []uint32
	// Project name to target
	Name []string
}

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

	if o.Setup {
		cmds = append(cmds, "-s")
	}

	if o.Clear {
		cmds = append(cmds, "-C")
	}

	if o.Check {
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

func (c *Command) Project(globalOpt GlobalOption, subopt ProjectCommandOption) error {
	c.SubOpt = subopt
	c.GlobalOpt = globalOpt
	return c.Execute()
}

func (c *Command) SetupProjectWithId(id uint32, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Setup: true,
		Path:  path,
		Depth: depth,
		Id:    []uint32{id},
	}

	gopt := GlobalOption{
		EnableExpertMode: true,
	}

	return c.Project(gopt, opt)
}

func (c *Command) SetupProjectWithName(name string, path string, depth uint32) error {
	opt := ProjectCommandOption{
		Setup: true,
		Path:  path,
		Depth: depth,
		Name:  []string{name},
	}

	gopt := GlobalOption{
		EnableExpertMode: true,
	}

	return c.Project(gopt, opt)
}

func (c *Command) ClearProjectWithId(id uint32, path string, depth uint32) error {
	opt := ProjectCommandOption{
    Clear: true,
		Path:  path,
		Depth: depth,
		Id:    []uint32{id},
	}

	gopt := GlobalOption{
		EnableExpertMode: true,
	}

	return c.Project(gopt, opt)
}

func (c *Command) ClearProjectWithName(name string, path string, depth uint32) error {
	opt := ProjectCommandOption{
    Clear: true,
		Path:  path,
		Depth: depth,
		Name:    []string{name},
	}

	gopt := GlobalOption{
		EnableExpertMode: true,
	}

	return c.Project(gopt, opt)
}

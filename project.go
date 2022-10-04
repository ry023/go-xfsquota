// +build linux
package xfsquota

import (
	"context"
	"fmt"
	"io"
	"regexp"
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
func (o projectCommandArgs) buildArgs() []string {
	args := []string{}
	args = append(args, "project")

	if o.opt.Depth != 0 {
		args = append(args, "-d")
		args = append(args, strconv.FormatUint(uint64(o.opt.Depth), 10))
	}

	if o.opt.Path != "" {
		args = append(args, "-p")
		args = append(args, o.opt.Path)
	}

	switch o.operation {
	case ProjectSetupOps:
		args = append(args, "-s")
	case ProjectClearOps:
		args = append(args, "-C")
	case ProjectCheckOps:
		args = append(args, "-c")
	}

	for _, id := range o.id {
		args = append(args, strconv.FormatUint(uint64(id), 10))
	}

	for _, name := range o.name {
		args = append(args, name)
	}

	return args
}

func (c *Command) OperateProjectWithId(ctx context.Context, op ProjectOpsType, id uint32, opt ProjectCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = &projectCommandArgs{
		id:        []uint32{id},
		operation: op,
		opt:       opt,
	}
	return c.Execute(ctx)
}

func (c *Command) SetupProjectWithId(ctx context.Context, id uint32, opt ProjectCommandOption) error {
	return c.OperateProjectWithId(ctx, ProjectSetupOps, id, opt)
}

func (c *Command) ClearProjectWithId(ctx context.Context, id uint32, opt ProjectCommandOption) error {
	return c.OperateProjectWithId(ctx, ProjectClearOps, id, opt)
}

func (c *Command) CheckProjectWithId(ctx context.Context, id uint32, opt ProjectCommandOption) error {
	err := c.OperateProjectWithId(ctx, ProjectClearOps, id, opt)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(c.systemStdoutBuf)
	if err != nil {
		return err
	}

	return c.checkOutput(b)
}

var ProjectCheckRegexp = regexp.MustCompile(`^(.*) - project inheritance flag is not set`)

type ProjectCheckError struct {
	Errors []error
}

func (e *ProjectCheckError) Error() string {
	if len(e.Errors) == 1 {
		return e.Errors[0].Error()
	} else {
		return fmt.Sprintf("%d error caused when checking a project.", len(e.Errors))
	}
}

var ProjectIdNotSetRegexp = regexp.MustCompile(`^(.*) - project identifier is not set`)

type ProjectIdNotSetError struct {
	Directory string
}

func (e *ProjectIdNotSetError) Error() string {
	return fmt.Sprintf("Project identifier is not set on directory (%s)", e.Directory)
}

var ProjectInheritanceFlagNotSetRegexp = regexp.MustCompile(`^(.*) - project inheritance flag is not set`)

type ProjectInheritanceFlagNotSetError struct {
	Directory string
}

func (e *ProjectInheritanceFlagNotSetError) Error() string {
	return fmt.Sprintf("Project inheritance flag is not set on directory (%s)", e.Directory)
}

func (c *Command) checkOutput(b []byte) error {
	var errs []error
	lines := strings.Split(string(b), "\n")
	for _, l := range lines {
		var submatches []string
		submatches = ProjectIdNotSetRegexp.FindStringSubmatch(l)
		if len(submatches) == 2 {
			errs = append(errs, &ProjectIdNotSetError{Directory: submatches[1]})
		}

		submatches = ProjectInheritanceFlagNotSetRegexp.FindStringSubmatch(l)
		if len(submatches) == 2 {
			errs = append(errs, &ProjectInheritanceFlagNotSetError{Directory: submatches[1]})
		}
	}

	if len(errs) > 0 {
		return &ProjectCheckError{Errors: errs}
	}
	return nil
}

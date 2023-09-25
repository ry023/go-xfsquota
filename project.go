package xfsquota

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ProjectCommander interface {
	// SetupDirectoryTree setup directory tree to project identified by projid.
	// Equeal to "project" subcommand with "-s" flag
	SetupDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error
	// ClearDirectoryTree clear directory tree to project identified by projid.
	// Equeal to "project" subcommand with "-C" flag
	ClearDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error
	// CheckDirectoryTree check directory tree to project identified by projid.
	// Equeal to "project" subcommand with "-c" flag
	CheckDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error
}

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
	ProjectDirTreeSetupOps = ProjectOpsType("Setup")
	ProjectDirTreeClearOps = ProjectOpsType("Clear")
	ProjectDirTreeCheckOps = ProjectOpsType("Check")
)

// Build 'project' subcommand
//
// format:
//
//	project [ -cCs [ -d depth ] [ -p path ] id | name ]
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
	case ProjectDirTreeSetupOps:
		args = append(args, "-s")
	case ProjectDirTreeClearOps:
		args = append(args, "-C")
	case ProjectDirTreeCheckOps:
		args = append(args, "-c")
	}

	for _, id := range o.id {
		args = append(args, strconv.FormatUint(uint64(id), 10))
	}

	args = append(args, o.name...)

	return args
}

func (c *Command) OperateDirectoryTree(ctx context.Context, op ProjectOpsType, id uint32, opt ProjectCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = &projectCommandArgs{
		id:        []uint32{id},
		operation: op,
		opt:       opt,
	}
	return c.Execute(ctx)
}

func (c *Command) SetupDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error {
	return c.OperateDirectoryTree(ctx, ProjectDirTreeSetupOps, projid, opt)
}

func (c *Command) ClearDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error {
	return c.OperateDirectoryTree(ctx, ProjectDirTreeClearOps, projid, opt)
}

func (c *Command) CheckDirectoryTree(ctx context.Context, projid uint32, opt ProjectCommandOption) error {
	err := c.OperateDirectoryTree(ctx, ProjectDirTreeCheckOps, projid, opt)
	if err != nil {
		return err
	}

	b, err := io.ReadAll(c.systemStdoutBuf)
	if err != nil {
		return err
	}

	return c.checkOutput(b)
}

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

var projectIdNotSetRegexp = regexp.MustCompile(`^(.*) - project identifier is not set`)

type ProjectIdNotSetError struct {
	Directory string
}

func (e *ProjectIdNotSetError) Error() string {
	return fmt.Sprintf("Project identifier is not set on directory (%s)", e.Directory)
}

var projectInheritanceFlagNotSetRegexp = regexp.MustCompile(`^(.*) - project inheritance flag is not set`)

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
		submatches = projectIdNotSetRegexp.FindStringSubmatch(l)
		if len(submatches) == 2 {
			errs = append(errs, &ProjectIdNotSetError{Directory: submatches[1]})
		}

		submatches = projectInheritanceFlagNotSetRegexp.FindStringSubmatch(l)
		if len(submatches) == 2 {
			errs = append(errs, &ProjectInheritanceFlagNotSetError{Directory: submatches[1]})
		}
	}

	if len(errs) > 0 {
		return &ProjectCheckError{Errors: errs}
	}
	return nil
}

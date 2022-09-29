package xfsquota

import (
	"bytes"
	"context"
	"io"
	"strings"
)

type Command struct {
	subCmdArgs subCommandArgs
	GlobalOpt  GlobalOption
	Stdout     io.Writer
	Stderr     io.Writer
	Binary     BinaryExecuter

	// The path argument can be used to specify mount points or device files which identify XFS filesystems. The output of the individual xfs_quota commands will then be restricted to the set of filesystems specified.
	// NOTE: This argument is optional on original xfs_quota but required on current go-xfsquota-wrapper version.
	FileSystemPath string

	systemStdoutBuf *bytes.Buffer
	systemStderrBuf *bytes.Buffer
}

type GlobalOption struct {
	// Equeal to "-p" flag on commandline.
	// Set the program name for prompts and some error messages, the default value is xfs_quota.
	ProgramName string
	// Equeal to "-x" flag on commandline.
	// Enable expert mode. All of the administrative commands which allow modifications to the quota system are available only in expert mode.
	EnableExpertMode bool
	// Equeal to "-d" flag on commandline.
	// Project names or numeric identifiers may be specified with this option, which restricts the output of the individual xfs_quota commands to the set of projects specified.
	Projects []string
}

// Interface to generate subcommands to specify for the -c option
type subCommandArgs interface {
	// Generate subcommand text
	buildArgs() []string
}

func NewCommand(binary BinaryExecuter, filesystemPath string, globalOpt *GlobalOption) *Command {
	cmd := &Command{
		Binary:         binary,
		FileSystemPath: filesystemPath,

		systemStdoutBuf: new(bytes.Buffer),
		systemStderrBuf: new(bytes.Buffer),
	}

	if globalOpt != nil {
		cmd.GlobalOpt = *globalOpt
	}

	return cmd
}

func (c *Command) Execute(ctx context.Context) error {
	args := c.buildArgs()

	var stdout io.Writer
	var stderr io.Writer

	if c.Stdout == nil {
		stdout = c.systemStdoutBuf
	} else {
		stdout = io.MultiWriter(c.Stdout, c.systemStdoutBuf)
	}

	if c.Stderr == nil {
		stderr = c.systemStderrBuf
	} else {
		stderr = io.MultiWriter(c.Stderr, c.systemStderrBuf)
	}

	return c.Binary.Execute(ctx, stdout, stderr, args...)
}

func (c *Command) buildArgs() []string {
	var args []string

	if c.GlobalOpt.EnableExpertMode {
		args = append(args, "-x")
	}

	if c.GlobalOpt.ProgramName != "" {
		args = append(args, "-p")
		args = append(args, c.GlobalOpt.ProgramName)
	}

	args = append(args, "-c")
	args = append(args, strings.Join(c.subCmdArgs.buildArgs(), " "))

	for _, d := range c.GlobalOpt.Projects {
		args = append(args, "-d")
		args = append(args, d)
	}

	args = append(args, c.FileSystemPath)
	return args
}

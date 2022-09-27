package xfsquota

import (
	"bytes"
	"fmt"
	"io"
)

type Command struct {
	SubOpt    SubCommandOption
	GlobalOpt GlobalOption
	Stdout    io.Writer
	Stderr    io.Writer
	Binary    BinaryExecuter

	systemStdoutBuf *bytes.Buffer
	systemStderrBuf *bytes.Buffer
}

func NewCommand(binary BinaryExecuter) *Command {
	return &Command{
		Binary: binary,
	}
}

func (c *Command) Execute() error {
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

	return c.Binary.Execute(stdout, stderr, args...)
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
	args = append(args, fmt.Sprintf("'%s'", c.SubOpt.SubCommandString()))

	for _, d := range c.GlobalOpt.Projects {
		args = append(args, "-d")
		args = append(args, d)
	}

	if c.GlobalOpt.Path != "" {
		args = append(args, c.GlobalOpt.Path)
	}
	return args
}

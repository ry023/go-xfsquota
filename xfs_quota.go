package xfsquota

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"

	v "github.com/hashicorp/go-version"
)

// xfs_quota wrapper client
type XfsQuotaClient struct {
	// The path to xfs_quota binary
	BinaryPath string
	// xfs_quota will only run if it satisfies the constraints of this version.
	VersionConstraint string
	// Ignore version checking if true. (Default is false)
	IgnoreVersionConstraint bool
	// Regexp used for parsing stdout of version command. (DefaultVersionCommandRegexp is used normally)
	VersionCommandRegexp *regexp.Regexp
}

var DefaultVersionCommandRegexp = regexp.MustCompile(`xfs_quota version\s(.*)\r?\n?$`)

type NewXfsQuotaClientOption func(*XfsQuotaClient) error

func NewXfsQuotaClient(binaryPath string, opts ...NewXfsQuotaClientOption) (*XfsQuotaClient, error) {
	c := &XfsQuotaClient{
		BinaryPath: binaryPath,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if err := c.validateBinary(); err != nil {
		return nil, err
	}

	return c, nil
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
	// The optional path argument can be used to specify mount points or device files which identify XFS filesystems. The output of the individual xfs_quota commands will then be restricted to the set of filesystems specified.
	Path string
}

type SubCommandOption interface {
	String() string
}

func (c *XfsQuotaClient) GetBinaryVersion() (string, error) {
	stdout, _, err := c.executeBinary("-V")
	if err != nil {
		return "", err
	}

	submatches := c.VersionCommandRegexp.FindSubmatch(stdout)
	if len(submatches) == 2 {
		return "", fmt.Errorf("Failed to parse version command stdout by c.VersionCommandRegexp(%s). (submatches=%v)", c.VersionCommandRegexp.String(), submatches)
	}

	version := string(submatches[1])
	return version, nil
}

func (c *XfsQuotaClient) validateBinary() error {
	if c.IgnoreVersionConstraint {
		return nil
	}

	constraints, err := v.NewConstraint(c.VersionConstraint)
	if err != nil {
		return err
	}

	// Check file existence
	if _, err = os.Stat(c.BinaryPath); err != nil {
		return err
	}

	version, err := c.GetBinaryVersion()
	if err != nil {
		return err
	}

	vv, err := v.NewVersion(version)
	if err != nil {
		return err
	}

	if !constraints.Check(vv) {
		return fmt.Errorf("xfs_quota binary violate version constraints! constraints=(%s), binary version=(%s)", constraints, version)
	}

	return nil
}

func (c *XfsQuotaClient) ExecuteCommand(commandOpt SubCommandOption, globalOpt GlobalOption) ([]byte, []byte, error) {
	var args []string

	if globalOpt.EnableExpertMode {
		args = append(args, "-x")
	}

	if globalOpt.ProgramName != "" {
		args = append(args, "-p")
		args = append(args, globalOpt.ProgramName)
	}

	args = append(args, "-c")
	args = append(args, fmt.Sprintf("'%s'", commandOpt.String()))

	for _, d := range globalOpt.Projects {
		args = append(args, "-d")
		args = append(args, d)
	}

	if globalOpt.Path != "" {
		args = append(args, globalOpt.Path)
	}

	return c.executeBinary(args...)
}

func (c *XfsQuotaClient) executeBinary(args ...string) ([]byte, []byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(c.BinaryPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, nil, err
	}

	stdoutBytes, err := io.ReadAll(&stdout)
	if err != nil {
		return nil, nil, err
	}

	stderrBytes, err := io.ReadAll(&stderr)
	if err != nil {
		return nil, nil, err
	}

	return stdoutBytes, stderrBytes, nil
}

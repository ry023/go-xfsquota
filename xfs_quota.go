package xfsquota

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	v "github.com/hashicorp/go-version"
)

// xfs_quota wrapper client
type XfsQuotaClient struct {
	// xfs_quota binary
	Binary BinaryExecuter
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
		Binary: &XfsQuotaBinary{
			Path: binaryPath,
		},
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

func (c *XfsQuotaClient) GetBinaryVersion() (string, error) {
	var stdout bytes.Buffer
	if err := c.Binary.Execute(&stdout, nil, "-V"); err != nil {
		return "", err
	}

	stdoutBytes, err := io.ReadAll(&stdout)
	if err != nil {
		return "", err
	}

	submatches := c.VersionCommandRegexp.FindSubmatch(stdoutBytes)
	if len(submatches) == 2 {
		return "", fmt.Errorf("Failed to parse version command stdout by c.VersionCommandRegexp(%s). (submatches=%v)", c.VersionCommandRegexp.String(), submatches)
	}

	version := string(submatches[1])
	return version, nil
}

func (c *XfsQuotaClient) validateBinary() error {
	if err := c.Binary.Validate(); err != nil {
		return err
	}

	if c.IgnoreVersionConstraint {
		return nil
	}

	constraints, err := v.NewConstraint(c.VersionConstraint)
	if err != nil {
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

func (c *XfsQuotaClient) Command(opt *GlobalOption) *Command {
  return NewCommand(c.Binary, opt)
}

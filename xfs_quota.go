package xfsquota

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"time"

	v "github.com/hashicorp/go-version"
)

// xfs_quota wrapper client.
type Client struct {
	// xfs_quota binary
	Binary BinaryExecuter
	// xfs_quota will only run if it satisfies the constraints of this version.
	VersionConstraint string
	// Ignore version checking if true. (Default is false)
	IgnoreVersionConstraint bool
	// Regexp used for parsing stdout of version command. (DefaultVersionCommandRegexp is used normally)
	VersionCommandRegexp *regexp.Regexp
}

const defaultVersionConstraint = ">= 5.13.0"

var defaultVersionCommandRegexp = regexp.MustCompile(`xfs_quota version\s(.*)\r?\n?$`)

func New(binaryPath string) (*Client, error) {
	c := &Client{
		Binary: &Binary{
			Path: binaryPath,
		},
		VersionConstraint:    defaultVersionConstraint,
		VersionCommandRegexp: defaultVersionCommandRegexp,
	}

	if err := c.validateBinary(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) GetBinaryVersion() (string, error) {
	stdout := new(bytes.Buffer)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	if err := c.Binary.Execute(ctx, stdout, nil, "-V"); err != nil {
		return "", err
	}

	submatches := c.VersionCommandRegexp.FindSubmatch(stdout.Bytes())
	if len(submatches) != 2 {
		return "", fmt.Errorf("Failed to parse version command stdout by c.VersionCommandRegexp(%s). (submatches=%v)", c.VersionCommandRegexp.String(), submatches)
	}

	version := string(submatches[1])
	return version, nil
}

func (c *Client) validateBinary() error {
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

func (c *Client) Command(filesystemPath string, opt *GlobalOption) Commander {
	return NewCommand(c.Binary, filesystemPath, opt)
}

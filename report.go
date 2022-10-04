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

type ReportCommander interface {
	Report(ctx context.Context, quotaType QuotaType, quotaTargetType QuotaTargetType, opt ReportCommandOption) (*ReportResult, error)
}

type reportCommandArgs struct {
	// Equal to `-gpu` flag on commandline.
	// Group/Project/User
	quotaType QuotaType
	// Equal to `-bir` flag on commandline.
	// Blocks/Inodes/Realtime
	quotaTargetType QuotaTargetType

	opt ReportCommandOption
}

type ReportCommandOption struct {
	// Equal to `-L` flag on commandline.
	// lower ID bounds to report on
	LowerId uint32
	// Equal to `-U` flag on commandline.
	// upper ID bounds to report on
	UpperId uint32
}

func (o reportCommandArgs) buildArgs() []string {
	args := []string{}
	args = append(args, "report")

	args = append(args, o.quotaType.Flag())
	args = append(args, o.quotaTargetType.Flag())

	// Force to numetric mode for parsing
	args = append(args, "-N")

	if o.opt.LowerId != 0 {
		args = append(args, "-L")
		args = append(args, strconv.FormatUint(uint64(o.opt.LowerId), 10))
	}

	if o.opt.UpperId != 0 {
		args = append(args, "-L")
		args = append(args, strconv.FormatUint(uint64(o.opt.UpperId), 10))
	}

	return args
}

func (c *Command) Report(ctx context.Context, quotaType QuotaType, quotaTargetType QuotaTargetType, opt ReportCommandOption) (*ReportResult, error) {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = reportCommandArgs{
		quotaType:       quotaType,
		quotaTargetType: quotaTargetType,
		opt:             opt,
	}
	err := c.Execute(ctx)
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(c.systemStdoutBuf)
	if err != nil {
		return nil, err
	}

	return parseHeadlessReportOutput(out, quotaType, quotaTargetType, c.FileSystemPath, "")
}

type ReportResult struct {
	ReportSets []ReportSet
}

type ReportSet struct {
	QuotaType       QuotaType
	QuotaTargetType QuotaTargetType
	MountPath       string
	DevicePath      string
	ReportValues    []ReportValue
}

type ReportValue struct {
	Id    uint32
	Used  uint32
	Soft  uint32
	Hard  uint32
	Grace uint32
}

func parseHeadlessReportOutput(out []byte, quotaType QuotaType, quotaTargetType QuotaTargetType, mountPath string, devicePath string) (*ReportResult, error) {
	outStr := string(out)

	// Trim head emptystring
	outStr = regexp.MustCompile(`^\n*`).ReplaceAllString(outStr, "")
	// Trim tail emptystring
	outStr = regexp.MustCompile(`\n*$`).ReplaceAllString(outStr, "")

	result := ReportResult{}

	reportSet := ReportSet{
		QuotaType:       quotaType,
		QuotaTargetType: quotaTargetType,
		MountPath:       mountPath,
		DevicePath:      devicePath,
	}

	// Split lines
	lines := strings.Split(outStr, "\n")
	for _, l := range lines {
		// line example:
		//   #100                12       1024      20480     05 [--------]

		// Trim spaces on line head
		l = regexp.MustCompile(`^\s*`).ReplaceAllString(l, "")

		v := regexp.MustCompile(`\s+`).Split(l, -1)
		if len(v) != 6 {
			return nil, fmt.Errorf("Failed to parse output")
		}

		id, err := strconv.Atoi(strings.TrimPrefix(v[0], "#"))
		if err != nil {
			return nil, err
		}

		used, err := strconv.Atoi(v[1])
		if err != nil {
			return nil, err
		}

		soft, err := strconv.Atoi(v[2])
		if err != nil {
			return nil, err
		}

		hard, err := strconv.Atoi(v[3])
		if err != nil {
			return nil, err
		}

		grace, err := strconv.Atoi(v[4])
		if err != nil {
			return nil, err
		}

		value := ReportValue{
			Id:    uint32(id),
			Used:  uint32(used),
			Soft:  uint32(soft),
			Hard:  uint32(hard),
			Grace: uint32(grace),
		}

		reportSet.ReportValues = append(reportSet.ReportValues, value)
	}
	result.ReportSets = append(result.ReportSets, reportSet)

	return &result, nil
}

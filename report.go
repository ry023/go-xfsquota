package xfsquota

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ReportCommandOption struct {
	// Equal to `-gpu` flag on commandline.
	// Group/Project/User
	QuotaTypes QuotaType
	// Equal to `-bir` flag on commandline.
	// Blocks/Inodes/Realtime
	QuotaTargetType QuotaTargetType
	// Equal to `-n` flag on commandline.
	// outputs the numeric ID instead of the name
	Numeric bool
	// Equal to `-L` flag on commandline.
	// lower ID bounds to report on
	LowerId uint32
	// Equal to `-U` flag on commandline.
	// upper ID bounds to report on
	UpperId uint32
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

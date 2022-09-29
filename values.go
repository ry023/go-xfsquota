// +build linux
package xfsquota

type QuotaType string

const (
	QuotaTypeGroup   = QuotaType("Group")
	QuotaTypeProject = QuotaType("Project")
	QuotaTypeUser    = QuotaType("User")
)

func (t QuotaType) Flag() string {
	switch t {
	case QuotaTypeGroup:
		return "-g"
	case QuotaTypeProject:
		return "-p"
	case QuotaTypeUser:
		return "-u"
	default:
		return ""
	}
}

type QuotaTargetType string

const (
	QuotaTargetTypeBlocks   = QuotaTargetType("Blocks")
	QuotaTargetTypeInodes   = QuotaTargetType("Inodes")
	QuotaTargetTypeRealtime = QuotaTargetType("Realtime")
)

func (t QuotaTargetType) Flag() string {
	switch t {
	case QuotaTargetTypeBlocks:
		return "-b"
	case QuotaTargetTypeInodes:
		return "-i"
	case QuotaTargetTypeRealtime:
		return "-r"
	default:
		return ""
	}
}

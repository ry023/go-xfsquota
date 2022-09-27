package xfsquota

type ReportCommandOption struct {
	// Equal to `-gpu` flag on commandline.
	// Group/Project/User
	QuotaType QuotaType
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

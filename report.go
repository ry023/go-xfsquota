package xfsquota

type ReportCommandOption struct {
	// Equal to `-g` flag on commandline.
	// Report quota of group type.
	Group bool
	// Equal to `-u` flag on commandline.
	// Report quota of user type.
	User bool
	// Equal to `-p` flag on commandline.
	// Report quota of project type.
	Project bool
	// Equal to `-b` flag on commandline.
	// Report quota of block type.
  Blocks bool
	// Equal to `-i` flag on commandline.
	// Report quota of inode type.
  Inodes bool
	// Equal to `-r` flag on commandline.
	// Report quota of realtime type.
  Realtime bool
	// Equal to `-a` flag on commandline.
  // Report all filesystems.
  AllFilesystems bool
  // Equal to `-n` flag on commandline.
  // outputs the numeric ID instead of the name
  Numeric bool
  // Equal to `-L` flag on commandline.
  // outputs the numeric ID instead of the name
  LowerID uint32
}


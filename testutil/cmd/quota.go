package cmd

type Quota struct {
	MountPath string
	Projects  map[uint32]Project
}

type Project struct {
	ID    uint32
	Paths []string
	// Block limits
	Bsoft uint32
	Bhard uint32
	// Inode limits
	Isoft uint32
	Ihard uint32
}

package xfsquota

import (
	"context"
	"fmt"
	"strconv"
)

type LimitCommander interface {
	LimitWithId(ctx context.Context, id uint32, quotaType QuotaType, opt LimitCommandOption) error
	LimitWithName(ctx context.Context, name string, quotaType QuotaType, opt LimitCommandOption) error
}

type LimitCommandOption struct {
	// Equal to `bsoft=N` (N > 0) argument on commandline.
	// Set quota block soft limits.
	Bsoft uint32
	// Equal to `bhard=N` (N > 0) argument on commandline.
	// Set quota block hard limits.
	Bhard uint32
	// Equal to `isoft=N` argument on commandline.
	// Set quota inode count soft limits.
	Isoft uint32
	// Equal to `ihard=N` (N > 0) argument on commandline.
	// Set quota inode count hard limits.
	Ihard uint32
	// Equal to `rtbsoft=N` (N > 0) argument on commandline.
	// Set quota realtime block soft limits.
	Rtbsoft uint32
	// Equal to `rtbhard=N` (N > 0) argument on commandline.
	// Set quota realtime block hard limits.
	Rtbhard uint32
	// Equal to `bsoft=0` (N > 0)  argument on commandline.
	// Reset quota block soft limits.
	ResetBsoft bool
	// Equal to `bhard=0` argument on commandline.
	// Reset quota block hard limits.
	ResetBhard bool
	// Equal to `isoft=0` argument on commandline.
	// Reset quota inode count soft limits.
	ResetIsoft bool
	// Equal to `ihard=0` argument on commandline.
	// Reset quota inode count hard limits.
	ResetIhard bool
	// Equal to `rtbsoft=0` argument on commandline.
	// Reset quota realtime block soft limits.
	ResetRtbsoft bool
	// Equal to `rtbhard=0` argument on commandline.
	// Reset quota realtime block hard limits.
	ResetRtbhard bool
}

type limitCommandArgs struct {
	// ID to target
	id []uint32
	// Name to target
	name []string
	// Equal to `-gpu` flag on commandline.
	// Group/Project/User
	quotaType QuotaType

	opt LimitCommandOption
}

func (o limitCommandArgs) buildArgs() []string {
	args := []string{}
	args = append(args, "limit")

	args = append(args, o.quotaType.Flag())

	if o.opt.ResetBhard {
		args = append(args, "bsoft=0")
	} else if o.opt.Bsoft != 0 {
		args = append(args, fmt.Sprintf("bsoft=%d", o.opt.Bsoft))
	}

	if o.opt.ResetBhard {
		args = append(args, "bhard=0")
	} else if o.opt.Bhard != 0 {
		args = append(args, fmt.Sprintf("bhard=%d", o.opt.Bhard))
	}

	if o.opt.ResetIsoft {
		args = append(args, "isoft=0")
	} else if o.opt.Isoft != 0 {
		args = append(args, fmt.Sprintf("isoft=%d", o.opt.Isoft))
	}

	if o.opt.ResetIhard {
		args = append(args, "ihard=0")
	} else if o.opt.Ihard != 0 {
		args = append(args, fmt.Sprintf("ihard=%d", o.opt.Ihard))
	}

	if o.opt.ResetRtbsoft {
		args = append(args, "rtbsoft=0")
	} else if o.opt.Rtbsoft != 0 {
		args = append(args, fmt.Sprintf("rtbsoft=%d", o.opt.Rtbsoft))
	}

	if o.opt.ResetRtbhard {
		args = append(args, "rtbhard=0")
	} else if o.opt.Rtbhard != 0 {
		args = append(args, fmt.Sprintf("rtbhard=%d", o.opt.Rtbhard))
	}

	for _, id := range o.id {
		args = append(args, strconv.FormatUint(uint64(id), 10))
	}

	args = append(args, o.name...)

	return args
}

func (c *Command) LimitWithId(ctx context.Context, id uint32, quotaType QuotaType, opt LimitCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = limitCommandArgs{
		id:        []uint32{id},
		quotaType: quotaType,
		opt:       opt,
	}
	return c.Execute(ctx)
}

func (c *Command) LimitWithName(ctx context.Context, name string, quotaType QuotaType, opt LimitCommandOption) error {
	c.GlobalOpt.EnableExpertMode = true // require expert mode
	c.subCmdArgs = limitCommandArgs{
		name:      []string{name},
		quotaType: quotaType,
		opt:       opt,
	}
	return c.Execute(ctx)
}

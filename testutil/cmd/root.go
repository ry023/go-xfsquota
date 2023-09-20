package cmd

import (
	"fmt"
	"os"

	"github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
)

var (
	version bool
	dummyX  bool
	dummyC  bool
)

var rootCmd = &cobra.Command{
	Use:   "fake_xfs_quota",
	Short: "Fake xfs_quota",
	Long:  `Fake xfs_quota.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println("fake_xfs_quota version 5.13.0")
			os.Exit(0)
		}
		if len(args) != 2 {
			return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", args)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c := args[0]
		m := args[1]
		cargs, err := shellwords.Parse(c)
		if err != nil {
			return err
		}
		switch cargs[0] {
		case "limit":
			return cmdLimit(os.Stdout, cargs[1:], m)
		case "report":
			return cmdReport(os.Stdout, cargs[1:], m)
		case "project":
			return cmdProject(os.Stdout, cargs[1:], m)
		default:
			return fmt.Errorf("fake_xfs_quota: invalid or unsupported command: %s", cargs[0])
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&dummyX, "x", "x", false, "Dummy flag")
	rootCmd.Flags().BoolVarP(&dummyC, "c", "c", false, "Dummy flag")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Show version and exit")
}

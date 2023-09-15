package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version bool

var rootCmd = &cobra.Command{
	Use:   "fake_xfs_quota",
	Short: "Fake xfs_quota",
	Long:  `Fake xfs_quota.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if version {
			fmt.Println("fake_xfs_quota version 5.13.0")
			os.Exit(0)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "Show version and exit")
}

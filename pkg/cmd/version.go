package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kube-resource",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kube-resource version v0.7.0")
	},
}

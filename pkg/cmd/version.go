package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version       string
	BuildDate     string
	GitCommitHash string
)

var versionTpl = `tools: kubectl-resource
 Version:           %v
 Go version:        %v
 Git commit:        %v
 Built:             %v
 OS/Arch:           %v
 Experimental:      false
 Repo: https://github.com/ysicing/kubectl-resource/releases/tag/%v
`

const (
	defaultVersion       = "0.0.0"
	defaultGitCommitHash = "a1b2c3d4"
	defaultBuildDate     = "Fri Jan 21 16:26:50 2022"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of kubectl-resource",
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			Version = defaultVersion
		}
		if BuildDate == "" {
			BuildDate = defaultBuildDate
		}
		if GitCommitHash == "" {
			GitCommitHash = defaultGitCommitHash
		}
		osarch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
		fmt.Printf(versionTpl, Version, runtime.Version(), GitCommitHash, BuildDate, osarch, Version)
	},
}

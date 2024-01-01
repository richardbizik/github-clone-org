package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(cloneCmd)
}

var RootCmd = &cobra.Command{
	Use:   "ghclone clone organization directory",
	Short: "ghclone is a cli tool to clone repositories of entire github organization",
	Long: `ghclone is a cli tool to clone repositories of entire github organization
	Usage: ghclone clone github.com/neovim`,
}

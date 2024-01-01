package main

import (
	"fmt"
	"os"

	"github.com/richardbizik/github-clone-org/internal/cmd"
)

func Execute() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

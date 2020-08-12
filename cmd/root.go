package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var (
	rootCmd = &cobra.Command{
		Use:   "rc",
		Short: "short description",
		Long:  "long description",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

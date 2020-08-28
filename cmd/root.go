package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// rootCmd 根命令
	rootCmd = &cobra.Command{
		Use:  "rc",
		Long: "remember and operate cli history command",
	}
)

// Execute 执行命令
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

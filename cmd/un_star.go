package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	// unStarCmd unStar命令
	unStarCmd = &cobra.Command{
		Use:   "us",
		Short: "unStar cmd",
		Long:  "unStar cmd by ID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {

				return errors.New("miss unStar id")
			}
			for _, v := range args {
				id, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("id %s invalid", v)
				}
				starIds[id] = struct{}{}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := handleUnStar()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

// init add the subcommand unStarCmd to rootCmd
func init() {
	rootCmd.AddCommand(unStarCmd)
}

// handleUnStar is not staring command
func handleUnStar() error {
	cache, err := getCommandS()
	if err != nil {
		return err
	}

	return starOrCancel(cache, starF)
}

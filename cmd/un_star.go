package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
)

var (
	// unStarIds 需要unStar的命令ID集合
	unStarIds = make(map[int]struct{})
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
				unStarIds[id] = struct{}{}
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := handleUnstar()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

// init 添加unStar子命令到根命令
func init() {
	rootCmd.AddCommand(unStarCmd)
}
func handleUnstar() error {
	cache, err := getCommands()
	if err != nil {
		return err
	}
	sortSlice := sortCommands(cache)
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := unStarIds[j]; exist {
			color.Info.Printf("unStar 【%s】\n", sortSlice[i].Cmd)
			cache[sortSlice[i].Cmd].Star = false
		}
		j++
	}
	// 保存
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, encode, 0666)
}

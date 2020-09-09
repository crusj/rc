package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
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
			err := handleUnStar()
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
func handleUnStar() error {
	cache, err := getCommandS()
	if err != nil {
		return err
	}
	var max uint32
	for _, c := range cache {
		if c.Times > max {
			max = c.Times
		}
	}
	sortSlice := CountSortS(cache, int(max))
	// record unstar cmds
	unStarCmds := make(map[string]struct{})
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := unStarIds[j]; exist {
			unStarCmds[sortSlice[i].Cmd] = struct{}{}
		}
		j++
	}
	// unset star
	for i, v := range cache {
		if _, exists := unStarCmds[v.Cmd]; exists {
			cache[i].Star = false
			color.Info.Printf("unstar 【%s】\n", cache[i].Cmd)
		}
	}
	// 保存
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, encode, 0666)
}

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
	// starIds 需要star的命令ID集合
	starIds = make(map[int]struct{})
	// starCmd star命令
	starCmd = &cobra.Command{
		Use:   "s",
		Short: "star cmd",
		Long:  "star cmd by ID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {

				return errors.New("miss star id")
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
			err := handleStar()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

// init 添加star子命令到根命令
func init() {
	rootCmd.AddCommand(starCmd)
}
func handleStar() error {
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
	// record star cmds
	starCmds := make(map[string]struct{})
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := starIds[j]; exist {
			starCmds[sortSlice[i].Cmd] = struct{}{}
		}
		j++
	}
	// set star
	for i, v := range cache {
		if _, exists := starCmds[v.Cmd]; exists {
			cache[i].Star = true
			color.Info.Printf("star 【%s】\n", cache[i].Cmd)
		}
	}
	// 保存
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, encode, 0666)
}

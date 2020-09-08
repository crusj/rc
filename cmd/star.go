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
	cache, err := getCommands()
	if err != nil {
		return err
	}
	sortSlice := sortCommands(cache)
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := starIds[j]; exist {
			color.Info.Printf("star 【%s】\n", sortSlice[i].Cmd)
			cache[sortSlice[i].Cmd].Star = true
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

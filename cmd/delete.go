/*
 * @Time : 2020/8/12 11:06 上午
 * @Author : 蒋龙
 * @File : add.go
 */
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
)

/**
 * 删除命令
 */

var (
	deleteIds = make(map[int]struct{}) // 需要删除的id
	deleteCmd = &cobra.Command{
		Use:   "d",
		Short: "add cmd",
		Long:  "add cmd from history file or update cmd frequency",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) <= 0 {
				return errors.New("miss delete id")
			}
			for _, v := range args {
				id, err := strconv.Atoi(v)
				if err != nil {
					return errors.New(fmt.Sprintf("id %s invalid", v))
				}
				deleteIds[id] = struct{}{}
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := handleDelete()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func handleDelete() error {
	cache, err := getCommands()
	if err != nil {
		return err
	}
	sortSlice := sortCommands(cache)
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := deleteIds[j]; exist {
			delete(cache, sortSlice[i].Cmd)
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

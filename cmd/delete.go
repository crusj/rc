package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gookit/color"

	"github.com/spf13/cobra"
)

/**
 * 删除命令
 */

var (
	// deleteIds 需要删除的命令ID集合
	deleteIds = make(map[int]struct{})
	// deleteCmd 删除命令
	deleteCmd = &cobra.Command{
		Use:   "d",
		Short: "delete cmd",
		Long:  "delete cmd by ID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {

				return errors.New("miss delete id")
			}
			for _, v := range args {
				id, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("id %s invalid", v)
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

// init 添加删除子命令到根命令
func init() {
	rootCmd.AddCommand(deleteCmd)
}

// handleDelete 执行删除命令
// 将命从缓存文件中删除
func handleDelete() error {
	cache, err := getCommands()
	if err != nil {
		return err
	}
	sortSlice := sortCommands(cache)
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := deleteIds[j]; exist {
			color.Red.Printf("delete 【%s】\n", sortSlice[i].Cmd)

			// restore tips
			restoreStr := bytes.Buffer{}
			restoreStr.WriteString(fmt.Sprintf(`restore use 【rc add "%s"`, sortSlice[i].Cmd))
			if len(sortSlice[i].Extra) != 0 {
				restoreStr.WriteString(fmt.Sprintf(` "%s"`, sortSlice[i].Extra))
			}
			if len(sortSlice[i].AliasID) != 0 {
				restoreStr.WriteString(fmt.Sprintf(` "%s"`, sortSlice[i].AliasID))
			}
			restoreStr.WriteString("】")
			color.Blue.Printf(restoreStr.String())

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

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

	// record cmds
	deleteCmd := make(map[string]struct{})
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := deleteIds[j]; exist {
			deleteCmd[sortSlice[i].Cmd] = struct{}{}
		}
		j++
	}
	// delete elem
	for i := 0; i < len(cache); i++ {
		if _, exists := deleteCmd[cache[i].Cmd]; exists {
			// delete info
			color.Red.Printf("delete 【%s】\n", cache[i].Cmd)
			// restore tips
			restoreStr := bytes.Buffer{}
			restoreStr.WriteString(fmt.Sprintf(`restore use 【rc add "%s"`, cache[i].Cmd))
			if len(cache[i].Extra) != 0 {
				restoreStr.WriteString(fmt.Sprintf(` "%s"`, cache[i].Extra))
			}
			if len(cache[i].AliasID) != 0 {
				restoreStr.WriteString(fmt.Sprintf(` "%s"`, cache[i].AliasID))
			}
			restoreStr.WriteString("】\n")
			color.Blue.Printf(restoreStr.String())

			cache = append(cache[:i], cache[i+1:]...)
			i--
		}
	}

	// 保存
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, encode, 0666)
}

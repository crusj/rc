package cmd

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"sort"
)

/**
 * 按照记录命令顺序频率进行排序
 */
type (
	SortSlice []*CacheDetail
)

func (s SortSlice) Len() int           { return len(s) }
func (s SortSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SortSlice) Less(i, j int) bool { return s[i].Times < s[j].Times }
func (s SortSlice) String() string {
	str := bytes.Buffer{}
	for i := len(s) - 1; i >= 0; i-- {
		str.WriteString(fmt.Sprintf("%d:%s\n", s[i].Times, s[i].Cmd))
	}
	return str.String()
}

var (
	listCmd = &cobra.Command{
		Use:   "l",
		Short: "short description",
		Long:  "long description",
		Run: func(cmd *cobra.Command, args []string) {
			err := handleList()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(listCmd)
}

func handleList() error {
	// 最近命令
	cache, err := getCommands()
	if err != nil {
		return err
	}
	s := sortCommands(cache)
	// 打印
	fmt.Println(s)
	return nil
}

// 命令排序
func sortCommands(cache Cache) SortSlice {
	s := make(SortSlice, len(cache))
	index := 0
	for _, v := range cache {
		s[index] = v
		index++
	}
	sort.Stable(s)
	return s
}

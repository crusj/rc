package cmd

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"sort"
	"strconv"
	"time"
)

/**
 * 按照记录命令顺序频率进行排序
 */
type (
	SortSlice []*CacheDetail
)

func (s SortSlice) Len() int      { return len(s) }
func (s SortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s SortSlice) Less(i, j int) bool {
	if s[i].Times == s[j].Times {
		ti, _ := time.Parse("2006-01-02 15:04:05", s[i].LastUpdate)
		tj, _ := time.Parse("2006-01-02 15:04:05", s[j].LastUpdate)
		return ti.Before(tj)
	}
	return s[i].Times < s[j].Times
}
func (s SortSlice) Render() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Id", "Cmd", "Extra", "FRE", "Last Update"})
	id := 1
	for i := len(s) - 1; i >= 0; i-- {
		table.Append([]string{
			strconv.Itoa(id),
			s[i].Cmd,
			s[i].Extra,
			strconv.FormatUint(uint64(s[i].Times), 10),
			s[i].LastUpdate,
		})
		id++
	}
	table.Render()
}

var (
	listCmd = &cobra.Command{
		Use:   "l",
		Short: "list commands",
		Long:  "list commands with table",
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
	s.Render()
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

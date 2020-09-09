package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type (
	// 已经排序的缓存(从小到大)
	SortSlice []*CacheDetail
)

func (s SortSlice) Len() int      { return len(s) }
func (s SortSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less 根据命令的频率进行排序，如果频率相同则根据时间先后进行排序
func (s SortSlice) Less(i, j int) bool {
	if s[i].Times == s[j].Times {
		ti, _ := time.Parse("2006-01-02 15:04:05", s[i].LastUpdate)
		tj, _ := time.Parse("2006-01-02 15:04:05", s[j].LastUpdate)

		return ti.Before(tj)
	}

	return s[i].Times < s[j].Times
}

// Render 以table的形式打印命令到终端
func (s SortSlice) Render() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"ID", "AID", "Extra", "Cmd", "FRE", "Last Update"})
	id := 1
	table.SetRowLine(rowLine)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i].Star {
			s[i].Extra = "✨ " + s[i].Extra
		}
		if star {
			fmt.Println(s[i].Star)
			if !s[i].Star {
				id++

				continue
			}
		}
		table.Append([]string{
			strconv.Itoa(id),
			s[i].AliasID,
			s[i].Extra,
			s[i].Cmd,
			strconv.FormatUint(uint64(s[i].Times), 10),
			s[i].LastUpdate,
		})
		id++
	}
	table.Render()
}

var (
	// rowLine is 是否显示row之间的行
	rowLine bool
	// star 是否只显示star的内容
	star    bool
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

// init 添加命令列表到子命令
func init() {
	listCmd.PersistentFlags().BoolVarP(&rowLine, "rowLine", "r", false, "table show row line")
	listCmd.PersistentFlags().BoolVarP(&star, "star", "s", false, "table only show star cmd")
	rootCmd.AddCommand(listCmd)
}

// handleList 执行显示命令列表
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

// sortCommands 将命令进行排序，返回已排序的命令
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

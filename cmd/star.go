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
		Short: "onlyShowStar cmd",
		Long:  "onlyShowStar cmd by ID",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {

				return errors.New("miss onlyShowStar id")
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

// starStatus represents the status of command
type starStatus bool

func (s starStatus) String() string {
	if s {
		return "onlyShowStar"
	}

	return "unStar"
}

var (
	starT starStatus = true
	starF starStatus = false
)

// init add subcommand starStatus to rootCmd
func init() {
	rootCmd.AddCommand(starCmd)
}

func handleStar() error {
	cache, err := getCommandS()
	if err != nil {
		return err
	}
	return starOrCancel(cache, starT)
}

// starOrCancel is staring or canceling command
func starOrCancel(cache CacheS, star starStatus) error {
	sortSlice := CountSortS(cache, int(cache.getMaxFre()))
	// record onlyShowStar cmds
	commands := make(map[string]struct{})
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if _, exist := starIds[j]; exist {
			commands[sortSlice[i].Cmd] = struct{}{}
		}
		j++
	}
	for i, v := range cache {
		if _, exists := commands[v.Cmd]; exists {
			cache[i].Star = bool(star)
			color.Info.Printf("%v 【%s】\n", star.String(), cache[i].Cmd)
		}
	}
	// 保存
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(cachePath, encode, 0666)
}

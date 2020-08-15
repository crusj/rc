/*
 * @Time : 2020/8/12 11:06 上午
 * @Author : 蒋龙
 * @File : add.go
 */
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

/**
 * 添加或更新命令频率
 */

type (
	CacheDetail struct {
		Cmd        string // 命令
		Times      uint32 // 记录频率
		LastUpdate string // 最后更新时间
		Extra      string // 备注
	}
	Cache map[string]*CacheDetail

	// 排序结构
	SortStruct struct {
		Times uint32
		Cmd   string
	}
)

var (
	addCmd      = &cobra.Command{
		Use:   "add",
		Short: "add cmd",
		Long:  "add cmd from history file or update cmd frequency",
		Run: func(cmd *cobra.Command, args []string) {
			err := handleAdd(append(args, "")[0])
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(addCmd)
}

func handleAdd(extra string) error {
	// 获取最后一条命令
	lastCommand := getLastCommand()
	// 获取已添加的命令
	cache, err := getCommands()
	if err != nil {
		return err
	}
	// 计数加一或添加信息的命令
	return setCommand(cache, lastCommand, extra)
}

// 获取最后一条命令
func getLastCommand() string {
	c1 := exec.Command("grep", "cmd:", historyPath)
	c2 := exec.Command("tail", "-n", "2")
	c3 := exec.Command("head", "-n", "1")
	c4 := exec.Command("awk", "-F", "cmd:", "{print $2}")
	c2.Stdin, _ = c1.StdoutPipe()
	c3.Stdin, _ = c2.StdoutPipe()
	c4.Stdin, _ = c3.StdoutPipe()
	stdOut := bytes.Buffer{}
	c4.Stdout = &stdOut
	_ = c4.Start()
	_ = c3.Start()
	_ = c2.Start()
	_ = c1.Run()
	_ = c2.Wait()
	_ = c3.Wait()
	_ = c4.Wait()
	return strings.TrimSpace(stdOut.String())
}

// 获取已添加的命令
func getCommands() (Cache, error) {
	file, err := os.OpenFile(cachePath, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var cache Cache
	if len(content) > 0 {
		err = json.Unmarshal(content, &cache)
		if err != nil {
			return nil, err
		}

		return cache, nil
	}
	return make(Cache), nil

}

// 添加命令
func setCommand(cache Cache, cmd, extra string) error {
	if _, exist := cache[cmd]; exist {
		cache[cmd].Times++
		cache[cmd].LastUpdate = time.Now().Format("2006-01-02 15:04:03")
		if len(extra) != 0 {
			cache[cmd].Extra = extra
		}
	} else {
		cache[cmd] = &CacheDetail{
			Cmd:        cmd,
			Times:      1,
			Extra:      extra,
			LastUpdate: time.Now().Format("2006-01-02 15:04:03"),
		}
	}
	encode, err := json.Marshal(cache)
	if err != nil {
		return err
	}
	// 文件不存在则创建文件
	if _, err = os.Stat(cachePath); os.IsNotExist(err) {
		if _, err := os.Create(cachePath); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(cachePath, encode, 0666)
}

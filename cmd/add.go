package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type (
	// CacheDetail 是存储一条命令的相关信息
	CacheDetail struct {
		Cmd        string // 命令
		AliasID    string // ID别名
		Times      uint32 // 记录频率
		LastUpdate string // 最后更新时间
		Extra      string // 备注
		Star       bool   // 是否star
	}
	// Cache 是所有已存储的命令
	Cache map[string]*CacheDetail
	// CacheS is cmd cache list
	CacheS []*CacheDetail

	//  SortStruct 是
	SortStruct struct {
		Times uint32
		Cmd   string
	}
)

var (
	// addCmd 添加命令历史命令最后一条到记录
	reCmd = &cobra.Command{
		Use:   "re",
		Short: "remember last history cmd",
		Long:  "re cmd from history file or update cmd frequency",
		Run: func(cmd *cobra.Command, args []string) {
			args = append(args, "", "")
			err := handleRe(args[0], args[1])
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	// addCmd 恢复命令,其实为跟具参数添加命令
	addCmd = &cobra.Command{
		Use:   "add",
		Short: "add cmd",
		Long:  "add cmd frequency",
		Run: func(cmd *cobra.Command, args []string) {
			args = append(args, "", "", "")
			err := handleAdd(args[0], args[1], args[2])
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

// init 添加添加命令到根命令
func init() {
	rootCmd.AddCommand(reCmd, addCmd)
}

// handleRe 执行添加命令
func handleRe(extra, alias string) error {
	// 获取最后一条命令
	lastCommand := getLastCommand()
	// 获取已添加的命令
	cache, err := getCommandS()
	if err != nil {
		return err
	}
	// 计数加一或添加信息的命令
	return setCommand(cache, lastCommand, extra, alias)
}

// handleAdd 执行恢复命令
func handleAdd(cmd, extra, alias string) error {
	// 获取已添加的命令
	cache, err := getCommandS()
	if err != nil {
		return err
	}
	// 计数加一或添加信息的命令
	return setCommand(cache, cmd, extra, alias)
}

// getLastCommand 返回用户历史命令中最后一条命令
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

// getCommandS returns slice of  command cache
func getCommandS() (CacheS, error) {
	file, err := os.OpenFile(cachePath, os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var cacheS CacheS
	if len(content) > 0 {
		err = json.Unmarshal(content, &cacheS)
		if err != nil {
			return nil, err
		}

		return cacheS, nil
	}

	return make(CacheS, 0), nil
}

// setCommand 添加新命令到缓存文件,如果命令已存在则增加命令的频率
func setCommand(cache CacheS, cmd, extra, alias string) error {
	if len(cmd) == 0 {
		return errors.New("命令不能为空")
	}
	exists := false
	for i, v := range cache {
		if v.Cmd == cmd {
			cache[i].Times++
			cache[i].LastUpdate = time.Now().Format("2006-01-02 15:04:03")
			if len(extra) != 0 {
				cache[i].Extra = extra
			}
			if len(alias) != 0 {
				cache[i].AliasID = alias
			}
			exists = true
			break
		}
	}
	if !exists {
		cache = append(cache, &CacheDetail{
			Cmd:        cmd,
			AliasID:    alias,
			Times:      1,
			Extra:      extra,
			LastUpdate: time.Now().Format("2006-01-02 15:04:03"),
		})
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

// cmdIncr 命令频率加一，如果命令不存在则新建命令
func cmdIncr(cmd string) error {
	cache, err := getCommandS()
	if err != nil {
		return err
	}

	return setCommand(cache, cmd, "", "")
}

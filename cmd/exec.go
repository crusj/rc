/*
 * @Time : 2020/8/15 1:44 下午
 * @Author : 蒋龙
 * @File : exec.go
 */
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// 执行命令历史中的命令
var (
	execId  = 0
	// 执行ID对应命令
	execCmd = &cobra.Command{
		Use:   "e",
		Short: "exec command",
		Long:  "exec ID Command",
		Args: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) <= 0 {
				return errors.New("miss exec id")
			}
			// 暂时支持一个命令执行
			execId, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := handleExec()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	// 编辑命令缓存文件
	cacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "vim cache",
		Long:  "vim cache",
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = handleEditCache()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
	// 赋值命令到clipboard
	cpCmd = &cobra.Command{
		Use:   "cp ID",
		Short: "copy command",
		Long:  "copy ID command to clipboard",
		Args: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) <= 0 {
				return errors.New("miss exec id")
			}
			// 暂时支持一个命令执行
			execId, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			err = handleCpExec()
			if err != nil {
				fmt.Println(err.Error())
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(execCmd, cacheCmd, cpCmd)
}
func handleEditCache() error {
	command := exec.Command("vim", cachePath)
	stdErr := &bytes.Buffer{}
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = stdErr
	if stdErr.Len() != 0 {
		fmt.Println(stdErr.String())
	}
	_ = command.Run()
	if stdErr.Len() != 0 {
		fmt.Println(stdErr.String())
	}

	return nil
}
func handleExec() error {
	execCmd, err := getCommand()
	if err != nil {
		return err
	}

	// 拆分命令
	execSplit := strings.Split(execCmd, " ")

	// 执行命令
	command := exec.Command(execSplit[0], execSplit[1:]...)
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	command.Stdout = stdOut
	cmdInfo := color.Green.Sprint("✔ " + execCmd + "\n")
	command.Stderr = stdErr
	cmdError := color.Red.Sprint("✗ " + execCmd + "\n")
	err = command.Run()
	if err != nil {
		fmt.Println(err.Error())
	}

	// 打印
	if stdErr.Len() != 0 { // 错误
		cmdError += stdErr.String()
		println(cmdError)
	} else { // 陈宫
		cmdInfo += stdOut.String()
		println(cmdInfo)
	}

	return nil
}
func handleCpExec() error {
	execCmd, err := getCommand()
	if err != nil {
		return err
	}

	return clipboard.WriteAll(execCmd)
}
func getCommand() (string, error) {
	var (
		execCmd string
	)

	// 获取执行命令
	cache, err := getCommands()
	if err != nil {
		return "", err
	}
	// 排序
	sortSlice := sortCommands(cache)
	// 排序为频率从大到小需要倒着取
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if j == execId {
			execCmd = sortSlice[i].Cmd
		}
		j++
	}

	if len(execCmd) == 0 {
		return "", errors.New("invalid command")
	}

	return execCmd, nil
}

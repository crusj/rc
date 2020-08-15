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
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"os/exec"
	"strconv"
	"strings"
)

// 执行命令历史中的命令
var (
	execId  = 0
	execCmd = &cobra.Command{
		Use:   "exec",
		Short: "exec cmd",
		Long:  "exec cmd by ID ",
		Args: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(args) < 0 {
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
)

func init() {
	rootCmd.AddCommand(execCmd)
}
func handleExec() error {
	var (
		err     error
		execCmd string
	)

	// 获取执行命令
	cache, err := getCommands()
	if err != nil {
		return err
	}
	sortSlice := sortCommands(cache)
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if j == execId {
			execCmd = sortSlice[i].Cmd
		}
		j++
	}
	if len(execCmd) == 0 {
		return errors.New("invalid command")
	}

	// 拆分命令
	execSplit := strings.Split(execCmd, " ")

	// 执行命令
	command := exec.Command(execSplit[0], execSplit[1:]...)
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	command.Stdout = stdOut
	cmdInfo := color.LightGreen.Sprint("✔ " + execCmd + "\n")
	command.Stderr = stdErr
	cmdError := color.Red.Sprint("✗ " + execCmd + "\n")

	_ = command.Run()

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

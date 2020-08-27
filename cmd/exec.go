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
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// 执行命令历史中的命令
var (
	// 命令ID
	execID = 0
	// 执行ID对应命令
	execCmd = &cobra.Command{
		Use:   "e",
		Short: "exec command",
		Long:  "exec ID Command",
		Args: func(cmd *cobra.Command, args []string) error {
			var err error
			if len(args) == 0 {
				return errors.New("miss exec id")
			}
			// 暂时支持一个命令执行
			execID, err = strconv.Atoi(args[0])
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
			err := handleEditCache()
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
			if len(args) == 0 {
				return errors.New("miss exec id")
			}
			// 暂时支持一个命令执行
			execID, err = strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := handleCpExec()
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
	commands := &Commands{}

	return commands.Handle(execCmd, "&&")
}

// 获取需要执行的命令
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
	// 根据排序ID获取需要执行的命令
	// ID排序频率从大到小需要倒着取
	j := 1
	for i := len(sortSlice) - 1; i >= 0; i-- {
		if j == execID {
			execCmd = sortSlice[i].Cmd
		}
		j++
	}

	if len(execCmd) == 0 {
		return "", errors.New("invalid command")
	}

	return execCmd, nil
}

// 复制命令到剪切板
func handleCpExec() error {
	execCmd, err := getCommand()
	if err != nil {
		return err
	}

	return clipboard.WriteAll(execCmd)
}

// 命令
type Commands struct {
	Command []string
	sep     string
}

// 根据分隔符拆分命令
func (re *Commands) split(cmds, sep string) {
	re.Command = strings.Split(cmds, sep)
	re.sep = sep
}
func (re *Commands) Handle(cmds, sep string) error {
	re.split(cmds, sep)
	switch sep {
	case "&&": // 连续执行多个命令
		for _, v := range re.Command {
			err := re.handleOne(v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (re *Commands) handleOne(cmd string) error {
	pieces := strings.Split(strings.TrimSpace(cmd), " ")
	// 执行命令
	command := exec.Command(pieces[0], pieces[1:]...)
	stdOut := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	command.Stdout = stdOut
	cmdInfo := color.Green.Sprint("✔ " + cmd + "\n")
	command.Stderr = stdErr
	cmdError := color.Red.Sprint("✗ " + cmd + "\n")
	err := command.Run()

	if err != nil {
		return errors.New(cmdError + color.Red.Sprintf(err.Error()))
	}

	// 打印
	if stdErr.Len() != 0 { // 错误
		cmdError += stdErr.String()
		return errors.New(cmdError)
	}
	cmdInfo += stdOut.String()
	println(cmdInfo)

	return nil
}

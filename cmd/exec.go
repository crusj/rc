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

	"github.com/atotto/clipboard"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var (
	// execID 需要执行的命令ID
	execID = 0
	// AliasID 需要执行命令的ID别名
	AliasID = ""
	// execCmd 执行子命令
	execCmd = &cobra.Command{
		Use:   "e",
		Short: "exec command",
		Long:  "exec ID Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("miss exec id")
			}
			// 匹配ID
			tmp, err := strconv.Atoi(args[0])
			if err != nil {
				AliasID = args[0]
			} else {
				execID = tmp
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			handleExec()
		},
	}
	// cacheCmd 编辑缓存子命令
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
	// cpCmd 复制命令子命令
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

// init 添加执行命令、编辑缓存、复制命令子命令到根命令
func init() {
	rootCmd.AddCommand(execCmd, cacheCmd, cpCmd)
}

// handleExec 执行子命令,并将命令频率加一
func handleExec() {
	var printOutput string
	execCmd, err := getCommand()
	if err != nil {
		fmt.Println(color.Red.Sprintf("获取执行命令失败:" + err.Error()))

		return
	}

	cmd := exec.Command(shell, "-c", execCmd)
	output := &bytes.Buffer{}
	stdErr := &bytes.Buffer{}
	cmd.Stdout = output
	cmd.Stderr = stdErr
	err = cmd.Run()
	if err != nil {
		printOutput = color.Red.Sprint("✗ "+execCmd+"\n") + stdErr.String() + err.Error()
	} else {
		printOutput = color.Green.Sprint("✔ "+execCmd+"\n") + output.String()
	}
	fmt.Println(printOutput)

	// 命令执行频率加一
	err = cmdIncr(execCmd)
	if err != nil {
		fmt.Println(color.Red.Sprint("FRE INCR FAILED: " + err.Error()))
	} else {
		fmt.Println(color.Green.Sprint("FRE INCR"))
	}
}

// getCommand 根据ExecId获取对应命令
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
		if execID > 0 && j == execID {
			execCmd = sortSlice[i].Cmd

			break
		}
		if len(AliasID) > 0 && sortSlice[i].AliasID == AliasID {
			execCmd = sortSlice[i].Cmd

			break
		}
		j++
	}

	if len(execCmd) == 0 {
		return "", errors.New("invalid command")
	}

	return execCmd, nil
}

// handleCpExec 复制命令到剪切板
func handleCpExec() error {
	execCmd, err := getCommand()
	if err != nil {
		return err
	}

	return clipboard.WriteAll(execCmd)
}

// handleEditCache 编辑缓存文件
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

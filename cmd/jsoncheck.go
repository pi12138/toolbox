//go:build jsoncheck

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Print bool
)

// jsoncheckCmd represents the jsoncheck command
var jsoncheckCmd = &cobra.Command{
	Use:   "jsoncheck <filename>",
	Short: "校验json文件和格式化打印json文件",
	Long:  `校验json文件和格式化打印json文件`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Printf("参数 <filename> 是必须的\n")
			return
		}
		filename := args[0]

		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("ReadFile error. %s\n", err)
			return
		}

		var a any
		if err := json.Unmarshal(data, &a); err != nil {
			fmt.Printf("json 文件内容存在问题. %s\n", err)
			return
		}

		if Print {
			if v, err := json.MarshalIndent(a, "", "  "); err != nil {
				fmt.Printf("JSON 格式化错误. %s\n", err)
				return
			} else {
				fmt.Printf("%s\n", string(v))
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(jsoncheckCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsoncheckCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsoncheckCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	jsoncheckCmd.Flags().BoolVarP(&Print, "print", "p", false, "print file data")
}

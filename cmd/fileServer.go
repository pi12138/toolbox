/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	folder string
	port   uint16
)

// fileServerCmd represents the fileServer command
var fileServerCmd = &cobra.Command{
	Use:   "fileServer",
	Short: "file server",
	Long:  `file server`,
	Run: func(cmd *cobra.Command, args []string) {
		// 指定文件服务器的根目录
		fs := http.FileServer(http.Dir(folder))

		// 设置路由
		http.Handle("/", fs)

		// 启动服务器并监听指定端口
		log(fmt.Sprintf("Server started on :%d\n", port))
		err := http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), nil)
		if err != nil {
			log(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(fileServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	fileServerCmd.Flags().StringVarP(&folder, "folder", "f", "tmp/", "folder path")
	fileServerCmd.Flags().Uint16VarP(&port, "port", "p", 7777, "folder path")
}

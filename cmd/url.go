/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

var (
	isDecode bool
	isEncode bool
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url <url-string>",
	Short: "url 编码和url解码工具",
	Long:  `url 编码和url解码工具`,
	Run: func(cmd *cobra.Command, args []string) {
		if isDecode {
			decodeUrl, err := url.QueryUnescape(args[0])
			if err != nil {
				fmt.Printf("decode %s error. %s", args[0], err)
				return
			}
			fmt.Println(decodeUrl)
		}
		if isEncode {
			encodeUrl := url.QueryEscape(args[0])
			fmt.Println(encodeUrl)
		}
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// urlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// urlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	urlCmd.Flags().BoolVarP(&isDecode, "decode", "d", false, "decode url")
	urlCmd.Flags().BoolVarP(&isEncode, "encode", "e", false, "encode url")
}

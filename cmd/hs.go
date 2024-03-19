/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	hs "github.com/pi12138/toolbox/internal/app/hospitalization_statistics"
	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	db "github.com/pi12138/toolbox/internal/pkg/database"
	"github.com/spf13/cobra"
)

// hsCmd represents the hs command
var hsCmd = &cobra.Command{
	Use:   "hs",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var kw string
		if len(args) == 0 {
			kw = "crawl"
		} else {
			kw = args[0]
		}

		switch kw {
		case "crawl":
			body, err := hs.Crawl()
			if err != nil {
				fmt.Printf("crawl error. %s\n", err)
			} else {
				fmt.Printf("crawl success. \n")
			}
			hs.SaveToJson(body)
		case "login":
			if err := hs.Login(); err != nil {
				fmt.Printf("login error. %s\n", err)
			} else {
				fmt.Printf("login success.\n")
			}
		case "migrate":
			model.Migrate(db.D())
		default:
			fmt.Printf("do nothing !!!\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(hsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

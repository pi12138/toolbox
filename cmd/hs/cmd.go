/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package hs

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

// hsCmd represents the hs command
var Cmd = &cobra.Command{
	Use:   "hs",
	Short: "hospitalization_statistics",
	Long:  `hospitalization_statistics`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("use -h show help")
	},
}

const (
	dateLayout = "2006-01-02"
)

func parseDate(s string) time.Time {
	date, err := time.Parse(dateLayout, s)
	if err != nil {
		log.Panicf("time.Parse error. %s\n", err)
	}
	return date
}

func defaultDate() string {
	return time.Now().Format(dateLayout)
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	Cmd.AddCommand(crawlCmd)
	Cmd.AddCommand(loginCmd)
	Cmd.AddCommand(migrateCmd)
	Cmd.AddCommand(statCmd)

}

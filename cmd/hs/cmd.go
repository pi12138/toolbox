/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package hs

import (
	"fmt"
	"log"
	"time"

	hs "github.com/pi12138/toolbox/internal/app/hospitalization_statistics"
	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	db "github.com/pi12138/toolbox/internal/pkg/database"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	UseJson bool
)

var (
	Start  string
	End    string
	Number string
)

// hsCmd represents the hs command
var Cmd = &cobra.Command{
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
			start := parseDate(Start)
			end := parseDate(End)

			dates := []string{}
			for t := start; t.Before(end) || t.Equal(end); t = t.AddDate(0, 0, 1) {
				dates = append(dates, t.Format(dateLayout))
			}
			for _, d := range dates {
				body, err := hs.Crawl(d, Number)
				prefix := fmt.Sprintf("date: %s, number: %s", d, Number)
				fmt.Printf("[%s] ", prefix)
				if err != nil {
					fmt.Printf("crawl error. %s.\n", err)
					return
				} else {
					fmt.Printf("crawl success. \n")
				}
				if UseJson {
					if err := hs.SaveToJson(body); err != nil {
						fmt.Printf("SaveToJson error. %s\n", err)
					}
				} else {
					if err := db.Transaction(func(tx *gorm.DB) error {
						return hs.SaveToDB(tx, body)
					}); err != nil {
						fmt.Printf("SaveToDB error. %s\n", err)
					}
				}
			}
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
	Cmd.Flags().BoolVarP(&UseJson, "use_json", "j", false, "save to json")
	Cmd.Flags().StringVarP(&Start, "start", "s", defaultDate(), "start date. (example: 2021-01-02)")
	Cmd.Flags().StringVarP(&End, "end", "e", defaultDate(), "end date.")
	Cmd.Flags().StringVarP(&Number, "number", "n", "", "query number")

}

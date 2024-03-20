package hs

import (
	"fmt"

	hs "github.com/pi12138/toolbox/internal/app/hospitalization_statistics"
	db "github.com/pi12138/toolbox/internal/pkg/database"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	UseJson bool

	Start  string
	End    string
	Number string
)

// hsCmd represents the hs command
var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "抓取数据",
	Long:  `抓取数据`,
	Run: func(cmd *cobra.Command, args []string) {

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

	},
}

func init() {

	crawlCmd.Flags().BoolVarP(&UseJson, "use_json", "j", false, "save to json")
	crawlCmd.Flags().StringVarP(&Start, "start", "s", defaultDate(), "start date. (example: 2021-01-02)")
	crawlCmd.Flags().StringVarP(&End, "end", "e", defaultDate(), "end date.")
	crawlCmd.Flags().StringVarP(&Number, "number", "n", "", "query number")
}

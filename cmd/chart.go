/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pi12138/toolbox/internal/pkg/chart"
	"github.com/spf13/cobra"
)

var (
	filename string
)

// chartCmd represents the chart command
var chartCmd = &cobra.Command{
	Use:   "chart",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("os.ReadFile %s error. %s", filename, err)
			return
		}

		var Data [][2]any
		if err := json.Unmarshal(data, &Data); err != nil {
			fmt.Printf("json.Unmarshal %s error. %s", data, err)
			return
		}

		chart.TimeSeries(Data)
		fs := http.FileServer(http.Dir("web/html"))
		log("running server at http://0.0.0.0:8089")
		log(http.ListenAndServe("0.0.0.0:8089", logRequest(fs)).Error())
	},
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func init() {
	rootCmd.AddCommand(chartCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chartCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chartCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	chartCmd.Flags().StringVarP(&filename, "filename", "f", "", "data filename")
}

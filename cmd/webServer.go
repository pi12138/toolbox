//go:build webServer

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var Port uint16

// webServerCmd represents the webServer command
var webServerCmd = &cobra.Command{
	Use:   "webServer",
	Short: "一个简单的 web server",
	Long:  `一个简单的 web server`,
	Run: func(cmd *cobra.Command, args []string) {
		Entry()
	},
}

func init() {
	rootCmd.AddCommand(webServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	webServerCmd.Flags().Uint16VarP(&Port, "port", "p", 8080, "server port")
}

func Entry() {
	http.HandleFunc("/", index)
	http.HandleFunc("/sleep", sleep)
	http.ListenAndServe(fmt.Sprintf(":%d", Port), nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "index")
}

func sleep(w http.ResponseWriter, r *http.Request) {
	sleepSecondsStr := r.URL.Query().Get("sleep")
	sleepSeconds, err := strconv.Atoi(sleepSecondsStr)
	if err != nil {
		http.Error(w, "Invalid sleep parameter", http.StatusBadRequest)
		return
	}

	time.Sleep(time.Duration(sleepSeconds) * time.Second)

	fmt.Fprintf(w, "Slept for %d seconds\n", sleepSeconds)
}

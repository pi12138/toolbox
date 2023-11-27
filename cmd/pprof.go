/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pi12138/toolbox/pkg/os/user"
	"github.com/spf13/cobra"
)

var DumpPath string
var Interval uint
var Forever bool

// const (
// 	DefaultInterval uint = uint(1 * time.Minute)
// )

var DefaultDumpPath string

var UriList = [3]string{
	"/debug/pprof/heap",
	"/debug/pprof/allocs",
	"/debug/pprof/goroutine",
	// "/debug/pprof/profile",
}

// pprofCmd represents the pprof command
var pprofCmd = &cobra.Command{
	Use:   "pprof <ip> <port>",
	Short: "定时抓取 pprof 的监控数据",
	Long:  `定时抓取 pprof 的监控数据.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println(args, DumpPath, Interval)
		if !Forever {
			for _, uri := range UriList {
				dumpFile(args[0], args[1], uri)
			}
			return
		}

		for {
			for _, uri := range UriList {
				dumpFile(args[0], args[1], uri)
			}
			time.Sleep(time.Duration(Interval * uint(time.Minute)))
		}
	},
}

func init() {
	rootCmd.AddCommand(pprofCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pprofCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pprofCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	homeDir, err := user.HomeDir()
	if err != nil {
		fmt.Printf("获取当前用户目录失败. %s\n", err)
		os.Exit(1)
	}
	DefaultDumpPath = filepath.Join(homeDir, "pprof")
	pprofCmd.Flags().StringVarP(&DumpPath, "dumppath", "d", DefaultDumpPath, "抓取数据文件保存位置")
	pprofCmd.Flags().UintVarP(&Interval, "interval", "i", 1, "抓取时间间隔(默认1分钟)")
	pprofCmd.Flags().BoolVarP(&Forever, "forever", "f", false, "是否持续运行")
}

func dumpFile(ip, port, path string) {
	Url := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%s", ip, port),
		Path:   path,
	}
	resp, err := http.Get(Url.String())
	if err != nil {
		log(fmt.Sprintf("Get %s error. %s", Url.String(), err))
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log("status code dont equal 200.")
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log(fmt.Sprintf("ReadAll error. %s", err))
		return
	}

	filename := toFilename(path)
	if err := os.WriteFile(filename, data, 0666); err != nil {
		log(fmt.Sprintf("WriteFile error. %s", err))
		return
	}
	log(fmt.Sprintf("dump file %s success.", filename))
}

func now() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func log(msg string) {
	fmt.Printf("[%s] %s\n", now(), msg)
}

// func logErr(msg string, err error) {
// 	fmt.Printf("[%s] %s. %s\n", now(), msg, err)
// }

func toFilename(path string) string {
	v := strings.Split(path, "/")
	suffix := v[len(v)-1]
	Now := now()
	body := strings.Replace(Now, " ", "-", 1)
	return filepath.Join(DumpPath, fmt.Sprintf("%s.%s", body, suffix))
}

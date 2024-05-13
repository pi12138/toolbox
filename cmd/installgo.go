/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pi12138/toolbox/internal/pkg/http"
	osPkg "github.com/pi12138/toolbox/pkg/os"
	"github.com/spf13/cobra"
)

var (
	InstallGoUrl        string
	InstallGoCacheFile  string
	InstallGoPathPrefix string
)

// installgoCmd represents the installgo command
var installgoCmd = &cobra.Command{
	Use:   "installgo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if InstallGoUrl == "" && InstallGoCacheFile == "" {
			fmt.Println("--file or --url must one.")
			os.Exit(1)
		}

		var filename string
		if InstallGoCacheFile != "" {
			if ok, err := osPkg.FileExist(InstallGoCacheFile); err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if !ok {
				fmt.Printf("%s don't exist.\n", InstallGoCacheFile)
				os.Exit(1)
			}
			filename = InstallGoCacheFile
		} else if InstallGoUrl != "" {
			var err error
			filename, err = http.ParseFilename(InstallGoUrl)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if err = http.Download(InstallGoUrl, filename); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("下载 %s 完毕.\n", filename)
		}

		var currentShell string
		var shellConfig string
		var home string
		home = os.Getenv("HOME")
		currentShell = os.Getenv("SHELL")
		switch currentShell {
		case "/bin/bash":
			shellConfig = filepath.Join(home, ".bashrc")
		case "/bin/zsh":
			shellConfig = filepath.Join(home, ".zshrc")
		default:
			fmt.Printf("暂不支持 %s\n", currentShell)
			os.Exit(1)
		}

		var input string
		GoRoot := filepath.Join(InstallGoPathPrefix, "go")
		GoBin := filepath.Join(GoRoot, "bin/go")
		GoPath := filepath.Join(home, ".go")
		BinPath := "/usr/bin/go"
		commands := [][]string{
			{"tar", "-zxvf", filename, "-C", "/tmp"},
			{"rm", "-rf", GoRoot},
			{"mv", "-f", "/tmp/go", GoRoot},
			{"ln", "-sf", GoBin, BinPath},
			{"echo", "export", "GOROOT=" + GoRoot, ">>", shellConfig},
			{"echo", "export", "GOPATH=" + GoPath, ">>", shellConfig},
		}

		fmt.Println("即将执行: ")
		for _, c := range commands {
			fmt.Println(exec.Command(c[0], c[1:]...))
		}
		fmt.Println("是否继续(输入yes继续): ")
		fmt.Scanf("%s", &input)
		fmt.Printf("用户输入: %s\n", input)
		if input != "yes" {
			fmt.Println("停止安装.")
			return
		}
		for _, c := range commands {
			if err := Exec(c[0], c[1:]...); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func Exec(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	fmt.Printf("exec %s\n", cmd.String())
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(installgoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installgoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installgoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	installgoCmd.Flags().StringVarP(&InstallGoUrl, "url", "u", "", "download go url")
	installgoCmd.Flags().StringVarP(&InstallGoCacheFile, "file", "f", "", "go archive package file (download from https://go.dev/dl/)")
	installgoCmd.Flags().StringVarP(&InstallGoPathPrefix, "prefix", "p", "/usr/local", "install prefix")

}

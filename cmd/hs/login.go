package hs

import (
	"fmt"

	hs "github.com/pi12138/toolbox/internal/app/hospitalization_statistics"
	"github.com/spf13/cobra"
)

// hsCmd represents the hs command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "获取用户身份",
	Long:  `获取登录数据`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := hs.Login(); err != nil {
			fmt.Printf("login error. %s\n", err)
		} else {
			fmt.Printf("login success.\n")
		}
	},
}

func init() {

}

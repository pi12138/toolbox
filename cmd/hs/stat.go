/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package hs

import (
	hs "github.com/pi12138/toolbox/internal/app/hospitalization_statistics"
	db "github.com/pi12138/toolbox/internal/pkg/database"
	"github.com/spf13/cobra"
)

// hsCmd represents the hs command
var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "统计数据",
	Long:  `统计数据`,
	Run: func(cmd *cobra.Command, args []string) {

		data := hs.Stat(db.D(), Start, End)
		// hs.ToTable(data)
		hs.ToExcel(data)
	},
}

func init() {
	statCmd.Flags().StringVarP(&Start, "start", "s", defaultDate(), "start date. (example: 2021-01-02)")
	statCmd.Flags().StringVarP(&End, "end", "e", defaultDate(), "end date.")
}

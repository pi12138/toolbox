/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package hs

import (
	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	db "github.com/pi12138/toolbox/internal/pkg/database"
	"github.com/spf13/cobra"
)

// hsCmd represents the hs command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate",
	Long:  `migrate`,
	Run: func(cmd *cobra.Command, args []string) {

		model.Migrate(db.D())
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

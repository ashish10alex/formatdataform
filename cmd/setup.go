/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates a supporting files need for formatting",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		yellow := color.New(color.FgYellow).SprintFunc()

		fmt.Println("Creating `.formatdataform` directory at the root of your project")
		os.Mkdir(".formatdataform", 0755)

        err := createFileFromText(pythonCode, ".formatdataform/sqlfluff_formatter.py")
		if err != nil {
            log.Println("Setup failed!!!")
			log.Fatalf(err.Error())
			return
		}

		err = createFileFromText(sqlfluffConfig, ".formatdataform/.sqlfluff")
		if err != nil {
            log.Println("Setup failed!!!")
			log.Fatalf(err.Error())
			return
		}

		fmt.Println("Setup complete")
		fmt.Println("Now you can run: ", yellow("`formatdataform format <path>`"))
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Creates a supporting files need for formatting",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

        yellow := color.New(color.FgYellow).SprintFunc()

		// create a directory .formatdataform
        fmt.Println("Creating", yellow(".formatdataform "), "directory at the root of your project")
		os.Mkdir(".formatdataform", 0755)
		//add the following lines of python code to the file .formatdataform/sqlfluff_formatter.py
		// add the following lines of python code to the file .formatdataform/sqlfluff_formatter.py
		f, err := os.Create(".formatdataform/sqlfluff_formatter.py")

        code := `
import sys
sqlfluff_config_path = sys.argv[1]
my_bad_query = sys.argv[2]

def fix_query(sqlfluff_config_path, my_bad_query):
    try:
        import sqlfluff
    except ImportError:
        return "sqlfluff is not installed"
    my_good_query = sqlfluff.fix(
                        my_bad_query,
                        dialect="bigquery",
                        config_path=str(sqlfluff_config_path),
                        )
    return my_good_query

print(fix_query(sqlfluff_config_path, my_bad_query)) # so that GoLang can read stdout
        `

		if err != nil {
			fmt.Println(err)
			return
		} else {
			f.WriteString(code)
            fmt.Println("sqlfluff_formatter.py file created at:", yellow(".formatdataform/sqlfluff_formatter.py"))
			f.Close()
            fmt.Println("Setup complete")
            fmt.Println("Now you can run: ", yellow("formatdataform format <path>"))
		}
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

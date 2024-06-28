/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "formats a file a directory depending on the next argument",
	Long: `formats a file a directory depending on the next argument
        To format a file: formatdataform format path/to/file.sqlx
        To format a directory: formatdataform format /dir/to/format
    `,
	Run: func(cmd *cobra.Command, args []string) {

		inplace, _ := cmd.Flags().GetBool("inplace")
		sqlfluffConfigPath := cmd.Flag("sqlfluff_config_path").Value.String()


        // make sure the .formatdataform directory exists if not create it
        os.Mkdir(".formatdataform", 0755)
		logFile, err := os.OpenFile(".formatdataform/formatdataform_logs.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer logFile.Close()

		logger := slog.New(slog.NewJSONHandler(logFile, nil))

		logger.Info("Formatting config",
			slog.String("sqlfluffConfigPath", sqlfluffConfigPath),
			slog.Bool("inplace", inplace),
		)


		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

        if sqlfluffConfigPath == ""{
            fmt.Printf(red("sqlfluff config file path is required \n"))
            return
        }

		if fileExists(sqlfluffConfigPath) == false {
			fmt.Printf(red("sqlfluff config file does not exist at: %v \n"),  sqlfluffConfigPath)
			return
		}

        if fileExists(".formatdataform/sqlfluff_formatter.py") == false {
            fmt.Print(yellow("sqlfluff_formatter.py file does not exist at: ", ".formatdataform/sqlfluff_formatter.py. Run: "))
            fmt.Printf("formatdataform setup \n")
            return
        }

        // If not file or directory path is supplied
		if len(args) == 0 {
            log.Fatalf(yellow(`No file or directory path supplied to command:  `, red(` formatdataform format <path>
                                                                                           ^^^^^`)))
            return
		}

        // If more than one file or directory path is supplied
        // TODO: Add support for multiple files or directories
        if len(args) > 1 {
            color.Set(color.FgYellow)
            fmt.Printf("Only supports one file or directory path at a time, you passed %v \n", len(args))
            color.Unset()
            return
        }

		fileOrDirPath := args[0]

		fileInfo, err := os.Stat(fileOrDirPath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
			return
		}

		if fileInfo.IsDir() {
			fmt.Println("Directory to format: ", green(fileOrDirPath))
			var sqlxFiles = findSqlxFiles(fileOrDirPath) // TODO: specify directory and depth to search for sql files here ?
			fmt.Println("Number of sqlx files found: ", green(len(*sqlxFiles))+"\n")

			if len(*sqlxFiles) == 0 {
				fmt.Println("No .sqlx files found in the directory: ", yellow(fileOrDirPath))
				return
			}

			var wg sync.WaitGroup
			for i := 0; i < len(*sqlxFiles); i++ {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					formatSqlxFile((*sqlxFiles)[i], inplace, sqlfluffConfigPath, logger)
				}(i)
			}
			wg.Wait()

		} else if !fileInfo.IsDir() {
			fmt.Println("File to format: ", green(fileOrDirPath))
            if filepath.Ext(fileOrDirPath) != ".sqlx" {
                fmt.Printf(red("Only .sqlx files are supported for formatting \n"))
                return
            }
			formatSqlxFile(fileOrDirPath, inplace, sqlfluffConfigPath, logger)
		} else {
			cmd.Help()
		}
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.Flags().StringP("file", "f", "", "file to format")
	formatCmd.Flags().StringP("dir", "d", "", "directory to format")
	formatCmd.Flags().BoolP("inplace", "i", true, "format the file in place")

}

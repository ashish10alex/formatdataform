/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
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

		green := color.New(color.FgGreen).SprintFunc()
		yellow := color.New(color.FgYellow).SprintFunc()
		red := color.New(color.FgRed).SprintFunc()

		inplace, _ := cmd.Flags().GetBool("inplace")
		sqlfluffConfigPath := cmd.Flag("sqlfluff_config_path").Value.String()

		fmt.Println("User args:\n\nSqlfluff config path: ", green(sqlfluffConfigPath))
		fmt.Println("Inplace: ", green(inplace))
		fmt.Print("\n")

        logger, logFile := setupLogger()
        defer logFile.Close()

		logger.Info("Formatting config passed by user",
			slog.String("sqlfluffConfigPath", sqlfluffConfigPath),
			slog.Bool("inplace", inplace),
		)

		if !inputArgsValid(args) {
			return
		}

		fileOrDirPath := args[0]
		valid, fileInfo := validFileOrDirPath(fileOrDirPath)
		if valid == false {
			return
		}

		setupValid, sqlfluffConfigPath := setupFilesAvailable(sqlfluffConfigPath)
		if setupValid == false {
			return
		}

		logger.Info("Formatting config post input validation",
			slog.String("sqlfluffConfigPath", sqlfluffConfigPath),
			slog.Bool("inplace", inplace),
            slog.Bool("isDir", fileInfo.IsDir()),
            slog.Bool("isFile", !fileInfo.IsDir()),
		)


		if fileInfo.IsDir() {
			fmt.Println("\nDirectory to format: ", green(fileOrDirPath))
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
			fmt.Println("\nFile to format: ", green(fileOrDirPath))
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

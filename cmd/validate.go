package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/fatih/color"
)


func inputArgsValid(args []string) bool {

	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// If not file or directory path is supplied
	if len(args) == 0 {
		log.Fatalf(yellow(`No file or directory path supplied to command:  `, red(` formatdataform format <path>
                                                                                           ^^^^^`)))
		return false
	}

	// If more than one file or directory path is supplied
	// TODO: Add support for multiple files or directories
	if len(args) > 1 {
		color.Set(color.FgYellow)
		fmt.Printf("Only supports one file or directory path at a time, you passed %v \n", len(args))
		color.Unset()
		return false
	}

	return true
}

func validFileOrDirPath(fileOrDirPath string) (bool, fs.FileInfo) {
	fileInfo, err := os.Stat(fileOrDirPath)
	if err != nil {
		log.Fatalf("Error opening file or directory path supplied: %v", err)
		return false, nil
	}
	return true, fileInfo
}


func setupFilesAvailable(sqlfluffConfigPath string) (bool, string) {

	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	if sqlfluffConfigPath == "" {
		sqlfluffConfigPath = ".formatdataform/.sqlfluff"
		fmt.Println("No sqlfluff config passed trying to using default ", sqlfluffConfigPath)
	}

	if fileExists(sqlfluffConfigPath) == false {
		fmt.Printf(yellow("Sqlfluff config file does not exist at: %v \n"), sqlfluffConfigPath)
		fmt.Println("Running ", green("`formatdataform setup` "), "to create a default config and supporting files")
		Setup()
	}

	if fileExists(".formatdataform/sqlfluff_formatter.py") == false {
		fmt.Print(yellow("sqlfluff_formatter.py file does not exist at: ", ".formatdataform/sqlfluff_formatter.py. Run: "))
		fmt.Printf("formatdataform setup \n")
		return false, sqlfluffConfigPath
	}

	return true, sqlfluffConfigPath
}

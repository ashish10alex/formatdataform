
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type sqlxFileMetaData struct {
	filepath        string
	numLines        int
	configStartLine int
	configEndLine   int
	configString    string
	queryString     string
	formattedQuery  string
}


func getSqlxFileMetaData(filepath string) (sqlxFileMetaData, error) {
	file, err := os.Open(filepath)
	if err != nil {
        fmt.Println("I am here")
		fmt.Println("Error opening file:", err)
		return sqlxFileMetaData{}, err
	}

	numLines, err := lineCounterV2(filepath)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return sqlxFileMetaData{}, err
	}

	// variables to keep track of where we are in the file
	var configStartLine = 0
	var configEndLine = 0
	var currentLine = 0
	var configString = ""
	var queryString = ""

	// flags to keep track of where we are in the file
	var isConfigBlock = false
	var isConfigBlockEnd = false
	var isInInnerConfigBlock = false
	var openCurlyBraceCount = 0
	var closeCurlyBraceCount = 0
	var queryBlockStarted = false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++
		var line = scanner.Text()

        //TODO: Check if this line is ever hit ?
		if err == io.EOF {
			break // End of file
		}

		// While we are in the config block, keep adding the lines to the configString
		if isConfigBlock == true && isConfigBlockEnd == false { // we are in the config block
			configString += line + "\n"
		}

		// If the line contains the word "config" and if we are not already in the config block, then we start the config block
		if strings.Contains(line, "config") && isConfigBlock == false {
			isConfigBlock = true
			configStartLine = currentLine
			configString += line + "\n"
		}

		// keep track of open and close curly braces while we are in the config block
		if strings.Contains(line, "{") && isConfigBlock == true {
			openCurlyBraceCount++
			if (openCurlyBraceCount != closeCurlyBraceCount) && (openCurlyBraceCount > 1) {
				isInInnerConfigBlock = true
			}
		}

		if strings.Contains(line, "}") && isConfigBlock == true { // TODO: breaks when you have } before the config block ends

			if configStartLine == 0 {
				configEndLine = 0
				// TODO: maybe we should return an error here
				fmt.Errorf("No config block found in file: %s", filepath)
			} else if isInInnerConfigBlock == true {
				closeCurlyBraceCount++
				isInInnerConfigBlock = false
			} else {
				configEndLine = currentLine
				isConfigBlockEnd = true
				isConfigBlock = false
			}
		}

		if isConfigBlockEnd == true && currentLine != configEndLine { // we are in the query block
			if line != "" {
				queryBlockStarted = true
			}
		}

		if isConfigBlockEnd == true && currentLine != configEndLine && queryBlockStarted { // we are in the query block
			if currentLine == numLines {
				queryString += line
			} else {
				queryString += line + "\n"
			}
		}

	}

	return sqlxFileMetaData{
		filepath:        filepath,
		numLines:        numLines,
		configStartLine: configStartLine,
		configEndLine:   configEndLine,
		configString:    configString,
		queryString:     queryString,
	}, nil
}

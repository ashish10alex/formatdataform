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
    preOperationsStartLine int
    preOperationsEndLine int
    preOperationsString string
	queryString     string
	formattedQuery  string
}

func getSqlxFileMetaData(filepath string) (sqlxFileMetaData, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return sqlxFileMetaData{}, err
	}

	numLines, err := countLinesInFile(filepath)

	if err != nil {
		fmt.Println("Error opening file:", err)
		return sqlxFileMetaData{}, err
	}

	// variables to keep track of where we are in the file
	var configStartLine = 0
	var configEndLine = 0
	var preOperationsStartLine = 0
	var preOperationsEndLine = 0
    var preOperationsString = ""
	var currentLineNumber = 0
	var configString = ""
	var queryString = ""

	// flags to keep track of where we are in the file
	var isConfigBlock = false
	var isConfigBlockEnd = false
	var isInInnerConfigBlock = false
	var openCurlyBraceCount = 0
	var closeCurlyBraceCount = 0
	var preOperationsBlockStarted = false
	var isInInnerPreOperationsBlock = false
    var isInPreOperationsBlock = false
	var queryBlockStarted = false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLineNumber++
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
			configStartLine = currentLineNumber
			configString += line + "\n"
		}

		// keep track of open and close curly braces while we are in the config block
		if strings.Contains(line, "{") && isConfigBlock == true {
			openCurlyBraceCount++
			if (openCurlyBraceCount != closeCurlyBraceCount) && (openCurlyBraceCount > 1) {
				isInInnerConfigBlock = true
			}
		}

		if strings.Contains(line, "}") && isConfigBlock == true { // TODO: breaks when you have curly brace before the config block ends

			if configStartLine == 0 {
				configEndLine = 0
				// TODO: maybe we should return an error here
				fmt.Errorf("No config block found in file: %s", filepath)
			} else if isInInnerConfigBlock == true {
				closeCurlyBraceCount++
				isInInnerConfigBlock = false // NOTE: does this mean that we only go to 1 nesting level ?
			} else {
				configEndLine = currentLineNumber
				isConfigBlockEnd = true
				isConfigBlock = false
			}
		}

		if isConfigBlockEnd == true && currentLineNumber != configEndLine { // query block started but looking for first non empty string
			if strings.Contains(line, "pre_operations") {
                isInPreOperationsBlock = true
				preOperationsBlockStarted = true
				openCurlyBraceCount = 0
				closeCurlyBraceCount = 0
				preOperationsStartLine = currentLineNumber
			}
		}

		if preOperationsBlockStarted == true && currentLineNumber != configEndLine {
            preOperationsString += line + "\n"
			if strings.Contains(line, "{") {
				openCurlyBraceCount++
				if (openCurlyBraceCount != closeCurlyBraceCount) && (openCurlyBraceCount > 1) {
					isInInnerPreOperationsBlock = true
				}
			}

			if strings.Contains(line, "}") {
				closeCurlyBraceCount++
				if isInInnerPreOperationsBlock == true {
					if closeCurlyBraceCount == openCurlyBraceCount {
                        isInPreOperationsBlock = false
						isInInnerPreOperationsBlock = false
						preOperationsEndLine = currentLineNumber
						preOperationsBlockStarted = false
					}
				}
			}
		}

		if isConfigBlockEnd == true && preOperationsBlockStarted == false && currentLineNumber != preOperationsEndLine && currentLineNumber != configEndLine {
			if line != "" {
				queryBlockStarted = true
			}
		}

		if queryBlockStarted && isInPreOperationsBlock == false && currentLineNumber != preOperationsEndLine { // in the query block
			if currentLineNumber == numLines {
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
        preOperationsStartLine: preOperationsStartLine,
        preOperationsEndLine: preOperationsEndLine,
        preOperationsString: preOperationsString,
		queryString:     queryString,
	}, nil
}

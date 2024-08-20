package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type ConfigBlockMeta struct {
	exsists            bool
	startOfConfigBlock int
	endOfConfigBlock   int
	configBlockContent string
}

type PreOpsBlockMeta struct {
	exsists                   bool
	startOfPreOperationsBlock int
	endOfPreOperationsBlock   int
	preOpsBlockContent        string
}

type PostOpsBlockMeta struct {
	exsists                    bool
	startOfpostOperationsBlock int
	endOfpostOperationsBlock   int
	postOpsBlockContent        string
}

type SqlBlockMeta struct {
	exsists                  bool
	startOfSqlBlock          int
	endOfSqlBlock            int
	sqlBlockContent          string
	formattedSqlBlockContent string
}

type sqlxParserMeta struct {
	filepath          string
	numLines          int
	configBlockMeta   ConfigBlockMeta
	preOpsBlocksMeta  []PreOpsBlockMeta
	postOpsBlocksMeta []PostOpsBlockMeta
	sqlBlocksMeta     SqlBlockMeta
}

func sqlxParser(filepath string) (sqlxParserMeta, error) {

	var inMajorBlock = false

	var startOfConfigBlock = 0
	var endOfConfigBlock = 0
	var configBlockExsists = false
	var configBlockContent = ""

	var preOpsBlocksMeta = []PreOpsBlockMeta{}
	var startOfPreOperationsBlock = 0
	var endOfPreOperationsBlock = 0

	var postOpsBlocksMeta = []PostOpsBlockMeta{}
	var startOfpostOperationsBlock = 0
	var endOfpostOperationsBlock = 0

	var startOfSqlBlock = 0
	var endOfSqlBlock = 0
	var sqlBlockExsists = false
	var sqlBlockContent = ""

	var isInInnerMajorBlock = false
	var innerMajorBlockCount = 0

	var currentBlock = ""
	var currentBlockContent = ""

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return sqlxParserMeta{}, err
	}

	i := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		var lineContents = scanner.Text() + "\n"

		if strings.Contains(lineContents, "config {") {
			inMajorBlock = true
			currentBlock = "config"
			startOfConfigBlock = i
			currentBlockContent += lineContents
		} else if strings.Contains(lineContents, "post_operations {") && !inMajorBlock {
			startOfpostOperationsBlock = i
			inMajorBlock = true
			currentBlock = "post_operations"
			currentBlockContent += lineContents
		} else if strings.Contains(lineContents, "pre_operations {") && !inMajorBlock {
			inMajorBlock = true
			currentBlock = "pre_operations"
			startOfPreOperationsBlock = i
			currentBlockContent += lineContents
		} else if strings.Contains(lineContents, "{") && inMajorBlock {
			if strings.Contains(lineContents, "}") {
				continue
			}
			isInInnerMajorBlock = true
			innerMajorBlockCount += 1
			currentBlockContent += lineContents
		} else if strings.Contains(lineContents, "}") && isInInnerMajorBlock && innerMajorBlockCount >= 1 && inMajorBlock {
			innerMajorBlockCount -= 1
			currentBlockContent += lineContents
		} else if strings.Contains(lineContents, "}") && innerMajorBlockCount == 0 && inMajorBlock {
			if currentBlock == "config" {
				currentBlockContent += lineContents
				configBlockContent = currentBlockContent
				endOfConfigBlock = i
				configBlockExsists = true
				currentBlock = ""
				currentBlockContent = ""
			} else if currentBlock == "pre_operations" {
				endOfPreOperationsBlock = i
				currentBlockContent += lineContents
				preOpsBlockMeta := PreOpsBlockMeta{
					exsists:                   true,
					startOfPreOperationsBlock: startOfPreOperationsBlock,
					endOfPreOperationsBlock:   endOfPreOperationsBlock,
					preOpsBlockContent:        currentBlockContent,
				}
				preOpsBlocksMeta = append(preOpsBlocksMeta, preOpsBlockMeta)
				currentBlock = ""
				currentBlockContent = ""
			} else if currentBlock == "post_operations" {
				endOfpostOperationsBlock = i
				currentBlockContent += lineContents
				postOpsBlockMeta := PostOpsBlockMeta{
					exsists:                    true,
					startOfpostOperationsBlock: startOfpostOperationsBlock,
					endOfpostOperationsBlock:   endOfpostOperationsBlock,
					postOpsBlockContent:        currentBlockContent,
				}
				postOpsBlocksMeta = append(postOpsBlocksMeta, postOpsBlockMeta)
				currentBlock = ""
				currentBlockContent = ""
			}
			inMajorBlock = false
		} else if strings.Contains(lineContents, "}") && isInInnerMajorBlock && innerMajorBlockCount >= 1 && !inMajorBlock {
			innerMajorBlockCount -= 1
			currentBlockContent += lineContents
		} else if lineContents != "\n" && !inMajorBlock {
			if startOfSqlBlock == 0 {
				startOfSqlBlock = i
				sqlBlockExsists = true
				sqlBlockContent += lineContents
				endOfSqlBlock = i
			} else {
				sqlBlockContent += lineContents
				endOfSqlBlock = i
			}
		} else if inMajorBlock {
			currentBlockContent += lineContents
		}
	}

	return sqlxParserMeta{
		filepath: filepath,
		numLines: i,
		configBlockMeta: ConfigBlockMeta{
			exsists:            configBlockExsists,
			startOfConfigBlock: startOfConfigBlock,
			endOfConfigBlock:   endOfConfigBlock,
			configBlockContent: configBlockContent,
		},
		preOpsBlocksMeta:  preOpsBlocksMeta,
		postOpsBlocksMeta: postOpsBlocksMeta,
		sqlBlocksMeta: SqlBlockMeta{
			exsists:                  sqlBlockExsists,
			startOfSqlBlock:          startOfSqlBlock,
			endOfSqlBlock:            endOfSqlBlock,
			sqlBlockContent:          sqlBlockContent,
			formattedSqlBlockContent: "",
		},
	}, nil
}

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

var ErrorFormattingSqlxFile = errors.New("Error formatting sqlx file")

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Walks the dataformRootDirectory and recursively finds sqlx files
func findSqlxFiles(dataformRootDirectory string) *[]string {
	var sqlxFilePath = filepath.Join(dataformRootDirectory)

	var sqlxFiles []string
	err := filepath.WalkDir(sqlxFilePath, func(path string, di fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".sqlx" {
			sqlxFiles = append(sqlxFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking the directory:", err)
	} else {
		return &sqlxFiles
	}
	return nil
}

func formatSqlCode(sqlxFileMetaData *sqlxParserMeta, pythonScriptPath string, sqlfluffConfigPath string, pythonExecutable string, logger *slog.Logger) error {
	queryString := *&sqlxFileMetaData.sqlBlocksMeta.sqlBlockContent

	cmd := exec.Command(pythonExecutable, pythonScriptPath, string(sqlfluffConfigPath), string(queryString))

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Error(stderr.String(), slog.String("file", sqlxFileMetaData.filepath), "error", err.Error())
        sqlxFileMetaData.sqlBlocksMeta.formattedSqlBlockContent = string(queryString)
		return ErrorFormattingSqlxFile
	}
	output := stdout.String()
	sql_fluff_not_installed := (strings.TrimSpace(output) == "sqlfluff is not installed")
	if sql_fluff_not_installed {
		log.Fatal(color.RedString("sqlfluff not installed. Please install sqlfluff using 'pip install sqlfluff'"))
	}
    sqlxFileMetaData.sqlBlocksMeta.formattedSqlBlockContent = output
	return nil
}

func finalFormmatedSqlxFileContents(sqlxFileMetaData *sqlxParserMeta) string {
    spaceBetweenBlocks := "\n\n"
    spaceBetweenSameOps := "\n"

    formattedQuery := ""

    preOpsBlocks := sqlxFileMetaData.preOpsBlocksMeta
    postOpsBlocks := sqlxFileMetaData.postOpsBlocksMeta

    preOpsBlockContent := ""
    if len(preOpsBlocks) > 0 {
        for _, preOpsBlock := range preOpsBlocks {
            preOpsBlockContent += preOpsBlock.preOpsBlockContent + spaceBetweenSameOps
        }
    }

    postOpsBlockContent := ""
    if len(postOpsBlocks) > 0 {
        for _, postOpsBlock := range postOpsBlocks {
            postOpsBlockContent += postOpsBlock.postOpsBlockContent + spaceBetweenSameOps
        }
    }

    formattedQuery = sqlxFileMetaData.configBlockMeta.configBlockContent +
                                                        spaceBetweenBlocks +
                                                        preOpsBlockContent +
                                                        spaceBetweenBlocks +
                                                        postOpsBlockContent +
                                                        spaceBetweenBlocks +
                                                        sqlxFileMetaData.sqlBlocksMeta.formattedSqlBlockContent
    return formattedQuery
}

func writeContentsToFile(sqlxFileMetaData *sqlxParserMeta, formattingError error) {

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

    filPathSeparator := string(os.PathSeparator)
    _definitions := "definitions" + filPathSeparator
	baseFilepath := strings.Split(sqlxFileMetaData.filepath, _definitions)
    formattedFilePath := filepath.Join("formatted", "definitions", baseFilepath[1])

	dirToCreate := formattedFilePath[:strings.LastIndex(formattedFilePath, filPathSeparator)]

	os.MkdirAll(dirToCreate, 0755) // TODO: make this configurable

    formattedQuery := finalFormmatedSqlxFileContents(sqlxFileMetaData)

	err := os.WriteFile(formattedFilePath, []byte(formattedQuery), 0664)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	switch formattingError {
	case nil:
		fmt.Printf("Formatted %d lines in %s\n", sqlxFileMetaData.numLines, yellow(formattedFilePath))
	case ErrorFormattingSqlxFile:
		fmt.Printf("Error formatting sqlx file: %s\n", red(sqlxFileMetaData.filepath))
	default:
	}
}

func writeContentsToFileInPlace(sqlxFileMetaData *sqlxParserMeta, formattingError error) {

	yellow := color.New(color.FgYellow).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()

    formattedQuery := finalFormmatedSqlxFileContents(sqlxFileMetaData)

	err := os.WriteFile(sqlxFileMetaData.filepath, []byte(formattedQuery), 0664)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	switch formattingError {
	case nil:
		fmt.Printf("Formatted %d lines in %s\n", sqlxFileMetaData.numLines, yellow(sqlxFileMetaData.filepath))
	case ErrorFormattingSqlxFile:
		fmt.Printf("Error formatting sqlx file: %s\n", red(sqlxFileMetaData.filepath))
	default:
	}
}

func formatSqlxFile(sqlxFilePath string, inplace bool, sqlfluffConfigPath string, pythonExecutable string, logger *slog.Logger) {
    sqlxFileMetaData, err := sqlxParser(sqlxFilePath)
	if err != nil {
		fmt.Println("Error finding config blocks:", err)
	} else {
        pythonScriptPath := filepath.Join(".formatdataform", "sqlfluff_formatter.py")
		formattingError := formatSqlCode(&sqlxFileMetaData, pythonScriptPath, sqlfluffConfigPath, pythonExecutable, logger)
		if inplace {
			writeContentsToFileInPlace(&sqlxFileMetaData, formattingError)
		} else {
			writeContentsToFile(&sqlxFileMetaData, formattingError)
		}

	}

}

func getIoReader(filepath string) (io.Reader, error) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	return file, nil
}

// Gives number of lines by reading the file in chunks, supposed to be faster than lineCounterV1 (https://stackoverflow.com/questions/24562942/golang-how-do-i-determine-the-number-of-lines-in-a-file-efficiently)

func lineCounterV3(reader io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := reader.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil
		case err != nil:
			return count, err
		}
	}
}

func countLinesInFile(filepath string) (int, error) {
	reader, err := getIoReader(filepath)
	if err != nil {
		return 0, err
	}
	return lineCounterV3(reader)
}

func createFileFromText(text string, filepath string) error {

	f, err := os.Create(filepath)

	if err != nil {
		return err
	} else {
		f.WriteString(text)
		fmt.Printf("file created at: `%s` \n",  filepath)
		f.Close()
	}
	return nil
}

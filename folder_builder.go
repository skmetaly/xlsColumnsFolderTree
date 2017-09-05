package main

import "fmt"
import (
	"bufio"
	"encoding/json"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
)

type Config struct {
	BasePath                 string `json:"basePath"`
	XlsPath                  string `json:"xlsPath"`
	FolderDepth              int    `json:"folderDepth"`
	StartProcessingFieldName string `json:"startProcessingFieldName"`
}

func processLevel(previousIndex int, previousPath string, level int, folderDepth int, dryRun bool, sheet *xlsx.Sheet) {

	rowCheck := false
	currentLevel := level + 1

	for currentIndex := previousIndex + 1; currentIndex < len(sheet.Rows) && rowCheck == false; currentIndex++ {
		if len(sheet.Rows[currentIndex].Cells) >= currentLevel+1 {

			if sheet.Rows[currentIndex].Cells[currentLevel-1].String() != "" {
				rowCheck = true
			} else {
				currentRowValue := sheet.Rows[currentIndex].Cells[currentLevel].String()
				if currentRowValue != "" {
					currentRowValue := previousPath + "/" + currentRowValue

					if _, err := os.Stat(currentRowValue); os.IsNotExist(err) {
						fmt.Println("Creating folder:" + currentRowValue)
						if dryRun == false {
							os.Mkdir(currentRowValue, os.ModePerm)
						}
					}

					if currentLevel < folderDepth {
						processLevel(currentIndex, currentRowValue, currentLevel, folderDepth, dryRun, sheet)
					}
				}
			}
		}
	}
}

func process(i int, depth int, basePath string, dryRun bool, sheet *xlsx.Sheet) {

	for cIndex := i + 1; cIndex < len(sheet.Rows); cIndex++ {

		if len(sheet.Rows[cIndex].Cells) != 0 && sheet.Rows[cIndex].Cells[0].String() != "" {
			firstRowValue := sheet.Rows[cIndex].Cells[0].String()

			if firstRowValue != "" {
				// Found first row. Try to see if it's already created
				firstRowPath := basePath + "/" + firstRowValue
				if _, err := os.Stat(firstRowPath); os.IsNotExist(err) {
					fmt.Println("Creating first level folder:" + firstRowPath)
					if dryRun == false {
						os.Mkdir(firstRowPath, os.ModePerm)
					}
				}

				processLevel(cIndex, firstRowPath, 0, depth, dryRun, sheet)
			}
		}
	}
}

func getConfig() Config {
	raw, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c Config
	json.Unmarshal(raw, &c)

	return c

}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func main() {

	c := getConfig()

	dryRun := false
	if len(os.Args) >= 2 {

		firstParam := os.Args[1]

		if firstParam == "test" {
			fmt.Println("Doing a dry run. No folders will be created")
			dryRun = true
		}
	}

	xlFile, err := xlsx.OpenFile(c.XlsPath)
	if err != nil {
		fmt.Println("Could not open file excel file " + c.XlsPath)
	}

	stat, err := exists(c.BasePath)

	if stat == false || err != nil {
		fmt.Println("Base folder doesn't exist or encountered an error accessing it")
	} else {

		for _, sheet := range xlFile.Sheets {
			for i, row := range sheet.Rows {

				if len(row.Cells) > 0 {
					text := row.Cells[0].String()
					if text == c.StartProcessingFieldName {
						process(i, c.FolderDepth, c.BasePath, dryRun, sheet)
					}
				}
			}
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter to finish")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
}

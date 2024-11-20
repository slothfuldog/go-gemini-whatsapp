package function

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var files *os.File

func init() {
	InitLogFileWin()
}

func InitLogFileWin() error {

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("Working Directory:", wd)

	// Define the filepath
	filePath := "data\\" // Using Windows filepath separator
	fileName := fmt.Sprintf("%sLIST_USER_%s.json", filePath, "JSON")

	// Print the resolved filepath for debugging
	resolvedFilePath := filepath.Join(wd, fileName)

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Join(wd, filePath), 0755)
	if err != nil {
		return err
	}

	// Open or create the file
	file, err := os.OpenFile(resolvedFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	files = file

	return nil
}

func CreateFile(detail string) {
	if files != nil {
		defer files.Sync() // Make sure logs are written before program exit

		wd, err := os.Getwd()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Working Directory:", wd)

		// Define the filepath
		filePath := "data\\" // Using Windows filepath separator
		fileName := fmt.Sprintf("%sLIST_USER_%s.json", filePath, "JSON")

		resolvedFilePath := filepath.Join(wd, fileName)

		err = os.WriteFile(resolvedFilePath, []byte(detail), 0644)
	}
}

func GetData() (string, error) {
	wd, err := os.Getwd()

	var content string

	if err != nil {
		return "", err
	}
	fmt.Println("Working Directory:", wd)

	// Define the filepath
	filePath := "data\\" // Using Windows filepath separator
	fileName := fmt.Sprintf("%sLIST_USER_%s.json", filePath, "JSON")

	resolvedFilePath := filepath.Join(wd, fileName)

	file, err := os.Open(resolvedFilePath)

	if err != nil {
		return "", err
	}

	scan := bufio.NewScanner(file)

	for scan.Scan() {
		content += scan.Text()
	}

	return content, nil
}

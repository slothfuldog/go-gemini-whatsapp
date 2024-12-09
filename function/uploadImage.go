package function

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func init() {
	iniUploadWin()
}

func iniUploadWin() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Define the filepath
	filePath := "image\\" // Using Windows filepath separator

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filepath.Join(wd, filePath), 0755)
	if err != nil {
		return err
	}

	// Open or create the file
	return nil
}

func UploadData(photo []byte) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	times := time.Now().Format("20060102150405")

	// Define the filepath
	filePath := "image\\" // Using Windows filepath separator
	fileName := fmt.Sprintf("%sPHOTO_%s%s", filePath, times, ".jpeg")

	// Print the resolved filepath for debugging
	resolvedFilePath := filepath.Join(wd, fileName)

	file, err := os.Create(resolvedFilePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return "", err
	}
	defer file.Close()

	reader := bytes.NewReader(photo)

	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Failed to save image:", err)
		return "", err
	}

	fmt.Println("Image saved to:", filePath)

	return fileName, nil
}

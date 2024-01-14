package fileutil

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

func CreateFolderIfNotExist(folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncFiles(sourceFolder, destinationFolder string) error {
	files, err := os.ReadDir(sourceFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		sourceFilePath := filepath.Join(sourceFolder, file.Name())
		destinationFilePath := filepath.Join(destinationFolder, file.Name())

		if _, err := os.Stat(destinationFilePath); os.IsNotExist(err) {
			if err := copyFile(sourceFilePath, destinationFilePath); err != nil {
				return err
			}
		}

		fmt.Println("File => ", file.Name())
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func applyTemplate(templateString string, data map[string]string) (string, error) {
	tmpl, err := template.New("local").Parse(templateString)
	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer
	if err := tmpl.Execute(&outputBuffer, data); err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}

func UpdateFilesInFolderWithTempaltes(folderPath string, templateValues map[string]string) error {
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	fmt.Println("UpdateFilesInFolderWithTempaltes =>", UpdateFilesInFolderWithTempaltes)
	for _, file := range files {
		filePath := filepath.Join(folderPath, file.Name())

		fmt.Println("filePath =>", filePath)
		fmt.Println("templateValues =>", templateValues)
		err := UpdateFileWithTempaltes(filePath, templateValues)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateFileWithTempaltes(filePath string, templateValues map[string]string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	newContent, err := applyTemplate(string(content), templateValues)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, []byte(newContent), os.ModePerm); err != nil {
		return err
	}
	return nil
}

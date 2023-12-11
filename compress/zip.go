package compress

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func Compress(target string) string {
	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	isDir := isDirectory(target)

	if isDir {
		err := addZipFiles(zipWriter, target, "")
		cobra.CheckErr(err)
	} else {
		fileName := filepath.Base(target)
		addZipFile(zipWriter, target, fileName)
	}

	err := zipWriter.Close()
	cobra.CheckErr(err)

	dumpDir, err := os.MkdirTemp("", ".datasync")
	cobra.CheckErr(err)

	zipFile := filepath.Join(dumpDir, "test.zip")

	file, err := os.Create(zipFile)
	cobra.CheckErr(err)
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	cobra.CheckErr(err)

	return zipFile
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	cobra.CheckErr(err)
	return fileInfo.IsDir()
}

func addZipFiles(writer *zip.Writer, basePath, pathInZip string) error {
	fileInfoArray, err := os.ReadDir(basePath)
	cobra.CheckErr(err)

	basePath = complementPath(basePath)
	pathInZip = complementPath(pathInZip)

	for _, fileInfo := range fileInfoArray {
		newBasePath := basePath + fileInfo.Name()
		newPathInZip := pathInZip + fileInfo.Name()

		if fileInfo.IsDir() {
			addDirectory(writer, newBasePath, newPathInZip)

			newBasePath = newBasePath + string(os.PathSeparator)
			newPathInZip = newPathInZip + string(os.PathSeparator)
			err = addZipFiles(writer, newBasePath, newPathInZip)
			cobra.CheckErr(err)
		} else {
			addZipFile(writer, newBasePath, newPathInZip)
		}
	}

	return nil
}

func addZipFile(writer *zip.Writer, targetFilePath, pathInZip string) {

	data, err := os.ReadFile(targetFilePath)
	cobra.CheckErr(err)

	fileInfo, err := os.Lstat(targetFilePath)
	cobra.CheckErr(err)

	header, err := zip.FileInfoHeader(fileInfo)
	cobra.CheckErr(err)

	header.Name = pathInZip
	header.Method = zip.Deflate
	w, err := writer.CreateHeader(header)
	cobra.CheckErr(err)

	_, err = w.Write(data)
	cobra.CheckErr(err)
}

func addDirectory(writer *zip.Writer, basePath string, pathInZip string) {
	fileInfo, err := os.Lstat(basePath)
	cobra.CheckErr(err)

	header, err := zip.FileInfoHeader(fileInfo)
	cobra.CheckErr(err)

	header.Name = pathInZip

	_, err = writer.CreateHeader(header)
	cobra.CheckErr(err)
}

func complementPath(path string) string {
	l := len(path)
	if l == 0 {
		return path
	}

	lastChar := path[l-1 : l]
	if lastChar == "/" || lastChar == "\\" {
		return path
	} else {
		return path + string(os.PathSeparator)
	}
}

func Decompress(dest, target string) {
	reader, err := zip.OpenReader(target)
	cobra.CheckErr(err)
	defer reader.Close()

	for _, zippedFile := range reader.File {
		path := filepath.Join(dest, zippedFile.Name)
		if zippedFile.FileInfo().IsDir() {
			err = os.MkdirAll(path, zippedFile.Mode())
			cobra.CheckErr(err)
		} else {
			createFileFromZipped(path, zippedFile)
		}
	}
}

func createFileFromZipped(path string, zippedFile *zip.File) {
	reader, err := zippedFile.Open()
	cobra.CheckErr(err)

	defer reader.Close()

	destFile, err := os.Create(path)
	cobra.CheckErr(err)

	defer destFile.Close()

	_, err = io.Copy(destFile, reader)
	cobra.CheckErr(err)
}

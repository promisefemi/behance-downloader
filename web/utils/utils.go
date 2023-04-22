package utils

import (
	"archive/zip"
	"fmt"
	"github.com/promisefemi/behance-downloader/core"
	"net/http"
	"path"
	"strings"
)

func HandleZipDownload(projectName string, images []string, rw http.ResponseWriter) {
	zipFileName := projectName + ".zip"

	zipFile := zip.NewWriter(rw)
	//defer zipFile.Close()
	for _, image := range images {
		fileName := path.Base(image)
		imageBody, _, _, err := core.ProcessDownload(image)
		if err != nil {
			continue
		}
		err = addFileToZip(zipFile, imageBody, fileName)
		if err != nil {
			fmt.Printf("Could not add %s to zip file -- %s", image, err)
			continue
		}
	}
	replaceZipFileName := strings.ReplaceAll(zipFileName, " ", "-")

	rw.Header().Set("Content-Disposition", "attachment; filename="+replaceZipFileName)
	rw.Header().Set("Content-Type", "application/zip")
	fmt.Printf("Downloading %s to client", replaceZipFileName)

	err := zipFile.Close()
	if err != nil {
		fmt.Printf("%s", err)
	}
	return
}

func addFileToZip(zipFile *zip.Writer, file []byte, fileName string) error {
	zipImage, err := zipFile.Create(fileName)
	if err != nil {
		return err
	}
	length, err := zipImage.Write(file)
	fmt.Printf("image write to zip length -- %d \n", length)
	if err != nil {
		return err
	}
	return nil
}
func DownloadImage(image string, rw http.ResponseWriter) {
	fileName := path.Base(image)
	fmt.Println(fileName)
	imageBody, contentType, contentLength, err := core.ProcessDownload(image)
	if err != nil {
		return
	}
	// buffer := bytes.NewBuffer(imageBody)
	fmt.Printf("%s --- %s --- %s", fileName, contentType, contentLength)
	rw.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Content-Length", contentLength)
	_, _ = rw.Write(imageBody)
	// io.Copy(rw, buffer)
}

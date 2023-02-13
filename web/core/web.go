package core

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/promisefemi/behance-downloader/core"
)

func Home(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		fmt.Println("Printing home page -- get")

		templ, err := parseTemplate("template/home.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(rw, nil)
		if err != nil {
			log.Fatal(err)
		}

		return
	}
	if r.Method == http.MethodPost {
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}

func Result(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("error parsing form")
		}

		url := r.PostForm.Get("url")

		behanceImage, err := core.ProcessLink(url)
		if err != nil {
			log.Fatal(err)
		}

		templ, err := parseTemplate("template/result.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(rw, behanceImage)
		if err != nil {
			log.Fatal(err)
		}
	}
	if r.Method == http.MethodGet {
		http.Redirect(rw, r, "/", http.StatusFound)
	}
}

func Download(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		selected := r.PostForm["selected[]"]

		if len(selected) > 1 {
			zipFileName:= "asdfasdfasdfsadfsadfasdf.zip"
			createFile, _ := os.Create(zipFileName)
			zipFile := zip.NewWriter(createFile)
			defer zipFile.Close()
			for _, image := range selected {
				// go func(image string) {
				//downlodImages(image, rw)
				// }(image)
				fileName := path.Base(image)

				imageBody, _, _, err := core.ProcessDownload(image)
				if err != nil {
					continue
				}
				addFileToZip(zipFile, imageBody, fileName)
			}



		} else {
			downlodImages(selected[0], rw)
		}

	}

}
func addFileToZip(zipFile *zip.Writer, file []byte, fileName string) error {

	zipImage, err := zipFile.Create(fileName)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(file)
	_, err = io.Copy(zipImage, reader)
	if err != nil {
		return err
	}

	return nil
}
func downlodImages(image string, rw http.ResponseWriter) {
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
	rw.Write(imageBody)
	// io.Copy(rw, buffer)
}

func parseTemplate(templateString string) (*template.Template, error) {
	fileName, err := filepath.Abs(templateString)
	if err != nil {
		return nil, err
	}
	templ, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

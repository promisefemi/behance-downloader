package core

import (
	"archive/zip"
	"fmt"
	"github.com/promisefemi/behance-downloader/core"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
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
		fileName := r.PostFormValue("projectName")
		if len(selected) > 1 {
			handleZipDownload(fileName, selected, rw)
		} else {
			downlodImages(selected[0], rw)
		}
		http.Redirect(rw, r, "", http.StatusFound)
	}

}
func handleZipDownload(fileName string, images []string, rw http.ResponseWriter) {
	zipFileName := fileName + ".zip"

	zipFile := zip.NewWriter(rw)
	defer zipFile.Close()
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

	rw.Header().Set("Content-Disposition", "attachment; filename="+zipFileName)
	rw.Header().Set("Content-Type", "application/zip")
	fmt.Printf("Downloading %s to client", zipFileName)

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

package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"

	"github.com/promisefemi/behance-downloader/core"
)

func Home(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		templ, err := parseTemplate("web/template/home.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(rw, nil)
		if err != nil {
			log.Fatal(err)

		}
		return
	}
}

func Result(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		url := r.PostForm.Get("url")

		behanceImage, err := core.ProcessLink(url)
		if err != nil {
			log.Fatal(err)
		}

		templ, err := parseTemplate("web/template/result.html")
		if err != nil {
			log.Fatal(err)
		}
		err = templ.Execute(rw, behanceImage)
		if err != nil {
			log.Fatal(err)
		}
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
			for _, image := range selected {
				// go func(image string) {
				downlodImages(image, rw)
				// }(image)
			}
		} else {
			downlodImages(selected[0], rw)
		}

	}

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

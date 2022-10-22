package web

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/promisefemi/behance-downloader/core"
)

func HandlePage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := homePage(rw, r)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.Method == http.MethodPost {
		err := getResult(rw, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func homePage(rw http.ResponseWriter, r *http.Request) error {
	templ, err := parseTemplate()
	if err != nil {
		return err
	}
	err = templ.Execute(rw, nil)
	if err != nil {
		return err
	}

	return nil
}

func getResult(rw http.ResponseWriter, r *http.Request) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}
	url := r.PostForm.Get("url")

	behanceImage, err := core.ProcessLink(url)
	if err != nil {
		return err
	}

	templ, err := parseTemplate()
	if err != nil {
		return err
	}
	err = templ.Execute(rw, behanceImage)
	if err != nil {
		return err
	}

	return nil
}

func parseTemplate() (*template.Template, error) {
	fileName, err := filepath.Abs("web/template/home.html")
	if err != nil {
		return nil, err
	}
	templ, err := template.ParseFiles(fileName)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

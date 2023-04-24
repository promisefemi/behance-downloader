package core

import (
	"embed"
	"fmt"
	"github.com/promisefemi/behance-downloader/core"
	"github.com/promisefemi/behance-downloader/web/utils"
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Directory embed.FS
}

func NewPageHandler(directory embed.FS) *Page {
	return &Page{
		directory,
	}
}

func (p Page) Home(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		fmt.Println("Printing home page -- get")

		templ, err := p.parseTemplate("template/home.html")
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

func (p Page) Result(rw http.ResponseWriter, r *http.Request) {

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

		//jsonByte, _ := json.MarshalIndent(behanceImage, " ", "   ")
		//fmt.Printf("%s", jsonByte)

		templ, err := p.parseTemplate("template/result.html")
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

func (p Page) Download(rw http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		selected := r.PostForm["selected[]"]
		fileName := r.PostFormValue("projectName")
		if len(selected) > 1 {
			utils.HandleZipDownload(fileName, selected, rw)
		} else {
			utils.DownloadImage(selected[0], rw)
		}
		http.Redirect(rw, r, "/", http.StatusSeeOther)
	}
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (p Page) parseTemplate(templateString string) (*template.Template, error) {

	templ, err := template.ParseFS(p.Directory, templateString)
	if err != nil {
		return nil, err
	}
	return templ, nil
}

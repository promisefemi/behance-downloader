package main

import (
	"embed"
	"fmt"
	"github.com/promisefemi/behance-downloader/web/core"
	"log"
	"net/http"
)

//go:embed template/*
var ViewsDirectory embed.FS

func main() {

	page := core.NewPageHandler(ViewsDirectory)

	//filesDirectory, err := ViewsDirectory.ReadFile("template/home.html")
	//if err != nil {
	//	fmt.Println(err)
	//}
	////var file []byte
	////fmt.Println(filesDirectory.Read(file))
	//fmt.Printf("\n%s\n", filesDirectory)

	http.HandleFunc("/", page.Home)
	http.HandleFunc("/result", page.Result)
	http.HandleFunc("/download", page.Download)
	http.HandleFunc("/favicon.ico", doNothing)

	fileServer := http.FileServer(http.FS(ViewsDirectory))

	http.Handle("/template/static/", http.StripPrefix("/", fileServer))

	fmt.Println("Server starting up @ http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func doNothing(w http.ResponseWriter, r *http.Request) {}

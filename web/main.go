package main

import (
	"embed"
	"fmt"
	"github.com/promisefemi/behance-downloader/web/core"
	"log"
	"net/http"
	"os"
)

//go:embed template/*
var ViewsDirectory embed.FS

func main() {

	page := core.NewPageHandler(ViewsDirectory)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	http.HandleFunc("/", page.Home)
	http.HandleFunc("/result", page.Result)
	http.HandleFunc("/download", page.Download)
	http.HandleFunc("/favicon.ico", doNothing)

	fileServer := http.FileServer(http.FS(ViewsDirectory))

	http.Handle("/template/static/", http.StripPrefix("/", fileServer))

	fmt.Printf("Server starting up @ port:%s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatal(err)
	}
}
func doNothing(w http.ResponseWriter, r *http.Request) {}

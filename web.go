package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/promisefemi/behance-downloader/web"
)

func main() {

	http.HandleFunc("/", web.Home)
	http.HandleFunc("/result", web.Result)
	http.HandleFunc("/download", web.Download)

	fileServer := http.FileServer(http.Dir("./web/template/static"))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	fmt.Println("Server starting up @ http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

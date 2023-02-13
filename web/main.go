package main

import (
	"fmt"
	"github.com/promisefemi/behance-downloader/web/core"
	"log"
	"net/http"

 )

func main() {

	http.HandleFunc("/", core.Home)
	http.HandleFunc("/result", core.Result)
	http.HandleFunc("/download", core.Download)
	http.HandleFunc("/favicon.ico", doNothing)

	fileServer := http.FileServer(http.Dir("./template/static"))

	http.Handle("/static/", http.StripPrefix("/static", fileServer))

	fmt.Println("Server starting up @ http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func doNothing(w http.ResponseWriter, r *http.Request){}

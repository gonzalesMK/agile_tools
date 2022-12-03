package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"
)

//go:embed build
var embeddedFiles embed.FS

func databases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	test := []string{}
	test = append(test, "Hello")
	test = append(test, "World")
	json.NewEncoder(w).Encode(test)
}

func main() {

	fsys, err := fs.Sub(embeddedFiles, "build")
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.FS(fsys)))
	http.Handle("/test", http.HandlerFunc(databases))

	http.ListenAndServe(":8050", nil)
}

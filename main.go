package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handled HTTP resuest
func (th *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	th.once.Do(func() {
		th.templ = template.Must(template.ParseFiles(filepath.Join("templates", th.filename)))
	})
	th.templ.Execute(w, r)
}
func main() {
	addr := flag.String("addr", ":8080", "The addr of the application")
	flag.Parse()

	r := NewRoom()
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)

	go r.run()

	log.Println("Starting web server on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}

}

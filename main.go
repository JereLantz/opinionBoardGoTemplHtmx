package main

import (
	"log"
	"net/http"
	"opinionBoardGoTemplHtmx/templates"

	"github.com/a-h/templ"
)

func handleRoot(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}
	templates.HelloWorld().Render(r.Context(),w)
}

func main() {
	handler := http.NewServeMux();
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}

	handler.HandleFunc("GET /", handleRoot)

	handler.Handle("GET /hello2", templ.Handler(templates.HelloWorld()))

	log.Printf("http server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

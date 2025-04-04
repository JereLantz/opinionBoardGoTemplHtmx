package main

import (
	"context"
	"log"
	"net/http"
	"opinionBoardGoTemplHtmx/templates"
)

func handleRoot(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
	templates.HelloWorld().Render(context.Background(), w)
	//w.Write([]byte("Hello world!"))
}

func main() {
	handler := http.NewServeMux();
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}

	handler.HandleFunc("GET /", handleRoot)

	log.Printf("http server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

package main

import (
	"database/sql"
	"log"
	"net/http"
	"opinionBoardGoTemplHtmx/templates/home"

	_ "github.com/mattn/go-sqlite3"
)

func handleRoot(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}

	home.Index().Render(r.Context(), w)
}

func initializeDbScheme(db *sql.DB) error{
	schemeInitQuery := `
	CREATE TABLE IF NOT EXISTS opinions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		opinion TEXT NOT NULL
	);
	`

	_, err := db.Exec(schemeInitQuery)
	if err != nil{
		return err
	}

	return nil
}
func connectToDB() (*sql.DB, error){
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil{
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	handler := http.NewServeMux();
	server := http.Server{
		Addr: ":42069",
		Handler: handler,
	}
	db, err := connectToDB()
	defer db.Close()
	if err != nil {
		log.Fatalf("error connecting to the database %s\n", err)
	}

	err = initializeDbScheme(db)
	if err != nil {
		log.Fatalf("could not initialize the database scheme %s\n",err)
	}

	handler.HandleFunc("GET /", handleRoot)

	log.Printf("http server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

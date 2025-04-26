package main

import (
	"database/sql"
	"log"
	"net/http"
	"opinionBoardGoTemplHtmx/templates/components"
	"opinionBoardGoTemplHtmx/templates/home"
	"opinionBoardGoTemplHtmx/utils"

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

func handleNewOpinion(db *sql.DB, w http.ResponseWriter, r *http.Request){
	var reRenderEmptyOpinion utils.Opinion

	errors,newOpinion, err := sanitizeInput(r)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if len(errors) > 0 {
		w.WriteHeader(400)
		components.ErrorDisplay(errors).Render(r.Context(),w)
		return
	}

	err = addOpinionDb(db, newOpinion)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	components.AddNewForm(reRenderEmptyOpinion).Render(r.Context(),w)
}

func addOpinionDb(db *sql.DB, newOpinion utils.Opinion) error{
	addOpinionQuery := `INSERT INTO opinions (title, opinion)
	values(?,?);`

	_, err := db.Exec(addOpinionQuery, newOpinion.Title, newOpinion.Opinion)
	if err != nil {
		return err
	}
	return nil
}

func sanitizeInput(r *http.Request) ([]string, utils.Opinion, error){
	var newOpinion utils.Opinion
	var errors []string

	err := r.ParseForm()
	if err != nil{
		return nil, newOpinion, err
	}

	title := r.FormValue("opinionTitle")
	opinionText := r.FormValue("opinion")

	if title == ""{
		errors = append(errors, "Please add a title.")
	}

	if opinionText == ""{
		errors = append(errors, "Please write down your opinion in text area.")
	}

	newOpinion.Title = title
	newOpinion.Opinion = opinionText

	return errors, newOpinion, nil
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

	handler.HandleFunc("POST /api/newopinion", func(w http.ResponseWriter, r *http.Request) {
		handleNewOpinion(db,w,r,)
	})

	log.Printf("http server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

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

func handleRoot(db *sql.DB, w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		w.WriteHeader(404)
		return
	}

	savedOpinions, err := fetchAllSavedOpinions(db)
	if err != nil {
		w.WriteHeader(503)
		return
	}

	home.Index(savedOpinions).Render(r.Context(), w)
}

func initializeDbScheme(db *sql.DB) error{
	schemeInitQuery := `
	CREATE TABLE IF NOT EXISTS opinions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		opinion TEXT NOT NULL,
		score INTEGER NOT NULL
	);
	`

	_, err := db.Exec(schemeInitQuery)
	if err != nil{
		return err
	}

	return nil
}

func fetchAllSavedOpinions(db *sql.DB) ([]utils.Opinion, error){
	var savedOpinions []utils.Opinion
	fetchAllQuery := `SELECT * FROM opinions
	ORDER BY id DESC;`

	row, err := db.Query(fetchAllQuery)
	if err != nil {
		return []utils.Opinion{}, err
	}

	for row.Next(){
		var newOpinion utils.Opinion

		err = row.Scan(&newOpinion.Id, &newOpinion.Title, &newOpinion.Opinion, &newOpinion.Score)
		if err != nil {
			return []utils.Opinion{}, err
		}

		savedOpinions = append(savedOpinions, newOpinion)
	}

	return savedOpinions, nil
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
	errors,newOpinion, err := sanitizeInput(r)
	invalidOpinion := utils.Opinion{
		Title : r.FormValue("opinionTitle"),
		Opinion: r.FormValue("opinion"),
	}

	if err != nil {
		w.Header().Add("HX-Reswap", "innerhtml")
		w.WriteHeader(400)
		components.ErrorResponse(errors, invalidOpinion).Render(r.Context(),w)
		return
	}

	if len(errors) > 0 {
		w.Header().Add("HX-Reswap", "outerHTML")
		w.WriteHeader(400)
		components.ErrorResponse(errors,invalidOpinion).Render(r.Context(),w)
		return
	}

	id, err := addOpinionDb(db, newOpinion)
	if err != nil {
		errors = append(errors, "Internal server error. Please try again later")
		w.Header().Add("HX-Reswap", "outerhtml")
		w.WriteHeader(500)
		components.ErrorDisplay(errors).Render(r.Context(),w)
		return
	}

	newOpinion.Id = id

	w.WriteHeader(200)
	components.Opinion(newOpinion).Render(r.Context(),w)
}

func addOpinionDb(db *sql.DB, newOpinion utils.Opinion) (int, error){
	addOpinionQuery := `INSERT INTO opinions (title, opinion, score)
	values(?, ?, 0);`

	result, err := db.Exec(addOpinionQuery, newOpinion.Title, newOpinion.Opinion)
	if err != nil {
		return -1, err
	}

	id, err := result.LastInsertId()

	return int(id), nil
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

	handler.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		handleRoot(db, w, r)
	})

	handler.HandleFunc("POST /api/newopinion", func(w http.ResponseWriter, r *http.Request) {
		handleNewOpinion(db,w,r)
	})

	handler.Handle("GET /index.js", http.FileServer(http.Dir("./")))

	log.Printf("http server started on port %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

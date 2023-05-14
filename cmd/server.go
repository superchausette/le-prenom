package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/superchausette/le-prenom/leprenom"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbName := flag.String("dbname", "", "Database to create or update")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *dbName == "" {
		flag.Usage()
		fmt.Println("Missing database name argument")
		return
	}

	fmt.Println("Opening sqlite database ", *dbName)
	db, err := gorm.Open(sqlite.Open(*dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&leprenom.FirstName{})
	db.AutoMigrate(&leprenom.Session{})
	db.AutoMigrate(&leprenom.SessionContent{})
	db.AutoMigrate(&leprenom.SessionNameStatus{})

	//DataStatDisplay(db)
	tmplFuncMap := template.FuncMap{
		"mod": func(value, modulo int) int {
			return value % modulo
		},
		"sessionTypeToStr": leprenom.SessionTypeToString}
	indexTmpl := template.Must(template.ParseFiles("template/index.html"))
	statsPartialTmpl := template.Must(template.ParseFiles("template/partial/stats.html"))
	listTmpl := template.Must(template.ParseFiles("template/list.html"))
	notFoundTmpl := template.Must(template.ParseFiles("template/404.html"))
	sessionListPartialTmpl := template.Must(template.New("SessionList").
		Funcs(tmplFuncMap).
		ParseFiles("template/partial/session_list.html"))
	firstNameListPartialTmpl := template.Must(template.New("List").
		Funcs(tmplFuncMap).
		ParseFiles("template/partial/firstname_list.html"))

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			notFoundTmpl.Execute(w, "")
			return
		}
		indexTmpl.Execute(w, "")
	}
	statsHandler := func(w http.ResponseWriter, r *http.Request) {
		statsPartialTmpl.Execute(w, leprenom.NewFirstNameStats(db))
	}
	listHandler := func(w http.ResponseWriter, r *http.Request) {
		// Check if this an htmx request
		htmx := r.Header.Get("HX-Request")
		if htmx != "" {
			listType := r.URL.Query().Get("type")
			first_names := leprenom.ListAllFirstName(db)
			err := firstNameListPartialTmpl.ExecuteTemplate(w, "firstname_list.html", first_names)
			if err != nil {
				fmt.Println("First Name List Partial Template error: ", err)
			}
			return
		}
		err = listTmpl.Execute(w, "")
	}
	newSessionHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		sessionName := r.Form.Get("session_name")
		sessionType := r.Form.Get("session_type")
		result := leprenom.CreateSession(sessionName, sessionType, db)
		w.Header().Set("HX-Trigger", "newSessionCreatedEvent")
		fmt.Fprintf(w, result)
	}
	listSessionHandler := func(w http.ResponseWriter, r *http.Request) {
		var sessions []leprenom.Session
		if err := db.Select("id", "name", "first_name_type").Find(&sessions).Error; err != nil {
			log.Fatal(err)
		}
		err = sessionListPartialTmpl.ExecuteTemplate(w, "session_list.html", sessions)
		if err != nil {
			fmt.Println("Session List Partial Template error: ", err)
		}
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/sessions/new", newSessionHandler)
	http.HandleFunc("/sessions/list", listSessionHandler)

	// Start the web server
	fmt.Println("Server is listening on http://localhost:9999/")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

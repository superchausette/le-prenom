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

	createSession := func(name string) string {
		fmt.Println("Creating session: ", name)
		result := db.Create(&leprenom.Session{Name: name})
		if result != nil && result.Error != nil {
			return fmt.Sprintf("Unable to create session '%s'", name)
		}
		return fmt.Sprintf("Succesfully created session '%s'", name)
	}

	//DataStatDisplay(db)

	indexTmpl := template.Must(template.ParseFiles("template/index.html"))
	notFoundTmpl := template.Must(template.ParseFiles("template/404.html"))
	sessionListPartialTmpl := template.Must(template.ParseFiles("template/partial/session_list.html"))
	statsPartialTmpl := template.Must(template.ParseFiles("template/partial/stats.html"))

	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			notFoundTmpl.Execute(w, "")
			return
		}
		indexTmpl.Execute(w, "")
	}
	newSessionHandler := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
		}
		result := createSession(r.Form.Get("session_name"))
		fmt.Fprintf(w, result)
	}
	listSessionHandler := func(w http.ResponseWriter, r *http.Request) {
		var sessions []leprenom.Session
		if err := db.Select("id", "name").Find(&sessions).Error; err != nil {
			log.Fatal(err)
		}
		sessionListPartialTmpl.Execute(w, sessions)
	}

	statsHandler := func(w http.ResponseWriter, r *http.Request) {
		statsPartialTmpl.Execute(w, leprenom.NewFirstNameStats(db))
	}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/sessions/new", newSessionHandler)
	http.HandleFunc("/sessions/list", listSessionHandler)
	http.HandleFunc("/stats", statsHandler)

	// Start the web server
	fmt.Println("Server is listening on http://localhost:9999/")
	log.Fatal(http.ListenAndServe(":9999", nil))
}

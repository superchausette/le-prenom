package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
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

	// Function used in templates
	tmplFuncMap := template.FuncMap{
		"mod": func(value, modulo int) int {
			return value % modulo
		},
		"sessionTypeToStr": leprenom.SessionTypeToString}
	templatesFiles := []string{"template/index.html",
		"template/list.html",
		"template/404.html",
		"template/partial/firstname_list.html",
		"template/partial/footer.html",
		"template/partial/header.html",
		"template/partial/session_list.html",
		"template/partial/stats.html",
	}
	templates := template.Must(template.New("templates").
		Funcs(tmplFuncMap).
		ParseFiles(templatesFiles...))
	rootHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			err := templates.ExecuteTemplate(w, "404.html", "")
			if err != nil {
				fmt.Println("Root Handler: unable to execute template", err)
			}
			return
		}
		err := templates.ExecuteTemplate(w, "index.html", "")
		if err != nil {
			fmt.Println("Root Handler: unable to execute template", err)
		}
	}
	statsHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		err := templates.ExecuteTemplate(w, "stats.html", leprenom.NewFirstNameStats(db))
		if err != nil {
			fmt.Println("Stats Handler: unable to execute template", err)
		}
	}
	listHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// Check if this an htmx request
		htmx := r.Header.Get("HX-Request")
		if htmx != "" {
			var firstNames []string
			listType := r.URL.Query().Get("type")
			switch listType {
			case "":
				firstNames = leprenom.ListAllFirstName(db)
			case "all":
				firstNames = leprenom.ListAllFirstName(db)
			case "boy":
				firstNames = leprenom.ListAllBoyFirstName(db)
			case "girl":
				firstNames = leprenom.ListAllGirlFirstName(db)
			case "unisex":
				firstNames = leprenom.ListAllUnisexFirstName(db)
			default:
				fmt.Println("ListHandler: unexpected type:", listType)
				return
			}
			err := templates.ExecuteTemplate(w, "firstname_list.html", firstNames)
			if err != nil {
				fmt.Println("First Name List Partial Template error: ", err)
			}
			return
		}
		err = templates.ExecuteTemplate(w, "list.html", "")
		if err != nil {
			fmt.Println("List Handler: unable to execute template", err)
		}
	}
	newSessionHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	listSessionHandler := func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var sessions []leprenom.Session
		if err := db.Select("id", "name", "first_name_type").Find(&sessions).Error; err != nil {
			log.Fatal(err)
		}
		err = templates.ExecuteTemplate(w, "session_list.html", sessions)
		if err != nil {
			fmt.Println("Session List Partial Template error: ", err)
		}
	}
	notFoundHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		err := templates.ExecuteTemplate(w, "404.html", "")
		if err != nil {
			fmt.Println("Root Handler: unable to execute template", err)
		}
	}

	router := httprouter.New()
	router.GET("/", rootHandler)
	router.GET("/stats", statsHandler)
	router.GET("/list", listHandler)
	router.POST("/sessions/new", newSessionHandler)
	router.GET("/sessions/list", listSessionHandler)
	router.NotFound = http.HandlerFunc(notFoundHandler)

	// Start the web server
	fmt.Println("Server is listening on http://localhost:9999/")
	log.Fatal(http.ListenAndServe(":9999", router))
}

package leprenom

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

type Server struct {
	DB       *gorm.DB
	Template *template.Template
}

func NewServer(db *gorm.DB) *Server {
	return &Server{DB: db, Template: NewTemplates()}
}

func (s *Server) RootHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		err := s.Template.ExecuteTemplate(w, "404.html", "")
		if err != nil {
			fmt.Println("Root Handler: unable to execute template", err)
		}
		return
	}
	err := s.Template.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		fmt.Println("Root Handler: unable to execute template", err)
	}
}

func (s *Server) StatsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := s.Template.ExecuteTemplate(w, "stats.html", NewFirstNameStats(s.DB))
	if err != nil {
		fmt.Println("Stats Handler: unable to execute template", err)
	}
}

func (s *Server) ListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check if this an htmx request
	htmx := r.Header.Get("HX-Request")
	if htmx != "" {
		var firstNames []string
		listType := r.URL.Query().Get("type")
		switch listType {
		case "":
			firstNames = ListAllFirstName(s.DB)
		case "all":
			firstNames = ListAllFirstName(s.DB)
		case "boy":
			firstNames = ListAllBoyFirstName(s.DB)
		case "girl":
			firstNames = ListAllGirlFirstName(s.DB)
		case "unisex":
			firstNames = ListAllUnisexFirstName(s.DB)
		default:
			fmt.Println("ListHandler: unexpected type:", listType)
			return
		}
		err := s.Template.ExecuteTemplate(w, "firstname_list.html", firstNames)
		if err != nil {
			fmt.Println("First Name List Partial Template error: ", err)
		}
		return
	}
	err := s.Template.ExecuteTemplate(w, "list.html", "")
	if err != nil {
		fmt.Println("List Handler: unable to execute template", err)
	}
}
func (s *Server) NewSessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}
	sessionName := r.Form.Get("session_name")
	sessionType := r.Form.Get("session_type")
	result := CreateSession(sessionName, sessionType, s.DB)
	w.Header().Set("HX-Trigger", "newSessionCreatedEvent")
	fmt.Fprintf(w, result)
}
func (s *Server) ListSessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var sessions []Session
	if err := s.DB.Select("id", "name", "first_name_type").Find(&sessions).Error; err != nil {
		log.Fatal(err)
	}
	err := s.Template.ExecuteTemplate(w, "session_list.html", sessions)
	if err != nil {
		fmt.Println("Session List Partial Template error: ", err)
	}
}
func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := s.Template.ExecuteTemplate(w, "404.html", "")
	if err != nil {
		fmt.Println("Root Handler: unable to execute template", err)
	}
}

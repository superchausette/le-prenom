package leprenom

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

func (s *Server) SessionHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sessionID, convErr := strconv.ParseUint(params.ByName("id"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Handler: Invalid session ID ", params.ByName("id"), " - ", convErr)
		http.Error(w, "SessionHandler: Invalid session ID", http.StatusBadRequest)
		return
	}
	var session Session
	sessionErr := s.DB.First(&session, sessionID).Error
	if sessionErr != nil {
		fmt.Println("Session Handler: Unable to get session ID ", sessionID, " - ", sessionErr)
		http.Error(w, fmt.Sprint("SessionHandler: Unable to get session ID", sessionID, sessionErr), http.StatusBadRequest)
		return
	}

	tmplParam := struct {
		ID        uint
		Name      string
		Type      string
		Total     uint
		Remaining uint
	}{uint(sessionID),
		session.Name,
		SessionTypeToString(session.FirstNameType),
		uint(TotalFirstNameAvailableForSession(session.FirstNameType, s.DB)),
		uint(RemainingFirstNameAvailableForSession(session.ID, session.FirstNameType, s.DB))}
	tmplErr := s.Template.ExecuteTemplate(w, "session.html", tmplParam)

	if tmplErr != nil {
		fmt.Println("SessionHandler Handler: unable to execute template", tmplErr)
	}
}

func (s *Server) SessionChoiceGetHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID, convErr := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Get Handler: Invalid session ID ", r.URL.Query().Get("id"), " - ", convErr)
		http.Error(w, "SessionGetChoiceGet: Invalid session ID", http.StatusBadRequest)
		return
	}
	tableIdx, convErr := strconv.ParseUint(r.URL.Query().Get("idx"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Get Handler: Invalid session ID ", r.URL.Query().Get("idx"), " - ", convErr)
		http.Error(w, "SessionGetChoiceGet: Invalid Index", http.StatusBadRequest)
		return
	}

	WriteSessionFirstNameEntry(uint(sessionID), uint(tableIdx), "SessionChoiceGetHandler", w, s.DB, s.Template)
}

func (s *Server) SessionChoiceKeepHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID, convErr := strconv.ParseUint(r.URL.Query().Get("sessionId"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Keep Handler: Invalid session ID ", r.URL.Query().Get("sessionId"), " - ", convErr)
		http.Error(w, "SessionGetChoiceKeep: Invalid session ID", http.StatusBadRequest)
		return
	}
	nameID, convErr := strconv.ParseUint(r.URL.Query().Get("nameId"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Keep Handler: Invalid session ID ", r.URL.Query().Get("nameId"), " - ", convErr)
		http.Error(w, "SessionGetChoiceKeep: Invalid name ID", http.StatusBadRequest)
		return
	}
	tableIdx, convErr := strconv.ParseUint(r.URL.Query().Get("tableIdx"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Keep Handler: Invalid session ID ", r.URL.Query().Get("idx"), " - ", convErr)
		http.Error(w, "SessionChoiceKeepHandler: Invalid Index", http.StatusBadRequest)
		return
	}

	// Write to db new value
	sessionContent := SessionContent{SessionID: uint(sessionID),
		FirstNameID: uint(nameID),
		Keep:        true,
	}
	createResult := s.DB.Create(&sessionContent)
	if createResult.Error != nil {
		fmt.Println("Session Choice Keep Handler: Unable to keep for session ",
			sessionID, " and name ", nameID, ":", createResult.Error)
		http.Error(w, "SessionGetChoiceKeep: Invalid name ID", http.StatusBadRequest)
	}
	w.Header().Set("HX-Trigger", "newFirstNameKept")
	// Return an new random first name to replace the one chosen
	WriteSessionFirstNameEntry(uint(sessionID), uint(tableIdx), "SessionChoiceKeepHandler", w, s.DB, s.Template)
}

func (s *Server) SessionChoiceRemoveHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID, convErr := strconv.ParseUint(r.URL.Query().Get("sessionId"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Remove Handler: Invalid session ID ", r.URL.Query().Get("sessionId"), " - ", convErr)
		http.Error(w, "SessionChoiceRemoveHandler: Invalid session ID", http.StatusBadRequest)
		return
	}
	nameID, convErr := strconv.ParseUint(r.URL.Query().Get("nameId"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Remove Handler: Invalid session ID ", r.URL.Query().Get("nameId"), " - ", convErr)
		http.Error(w, "SessionChoiceRemoveHandler: Invalid name ID", http.StatusBadRequest)
		return
	}
	tableIdx, convErr := strconv.ParseUint(r.URL.Query().Get("tableIdx"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice Remove Handler: Invalid session ID ", r.URL.Query().Get("idx"), " - ", convErr)
		http.Error(w, "SessionChoiceRemoveHandler: Invalid Index", http.StatusBadRequest)
		return
	}

	// Write to db new value
	sessionContent := SessionContent{SessionID: uint(sessionID),
		FirstNameID: uint(nameID),
		Keep:        false,
	}
	createResult := s.DB.Create(&sessionContent)
	if createResult.Error != nil {
		fmt.Println("Session Choice remove Handler: Unable to remove for session ",
			sessionID, " and name ", nameID, ":", createResult.Error)
		http.Error(w, "SessionChoiceRemoveHandler: Invalid name ID", http.StatusBadRequest)
	}

	// Return an new random first name to replace the one chosen
	WriteSessionFirstNameEntry(uint(sessionID), uint(tableIdx), "SessionChoiceKeepHandler", w, s.DB, s.Template)
}

func (s *Server) SessionChoiceListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionID, convErr := strconv.ParseUint(r.URL.Query().Get("sessionId"), 10, 64)
	if convErr != nil {
		// Handle the error if the conversion fails
		fmt.Println("Session Choice List Handler: Invalid session ID ", r.URL.Query().Get("sessionId"), " - ", convErr)
		http.Error(w, "SessionChoiceListHandler: Invalid session ID", http.StatusBadRequest)
		return
	}
	// Retrieve all first name kept
	firstNames := SessionFirstNameKept(uint(sessionID), s.DB)

	tmplErr := s.Template.ExecuteTemplate(w, "session_first_name_kept.html", firstNames)
	if tmplErr != nil {
		fmt.Println("SessionChoiceListHandler: unable to execute template", tmplErr)
	}
}

func (s *Server) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	err := s.Template.ExecuteTemplate(w, "404.html", "")
	if err != nil {
		fmt.Println("Root Handler: unable to execute template", err)
	}
}

func (s *Server) SetupRoutes(router *httprouter.Router) {
	router.GET("/", s.RootHandler)
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/stats", s.StatsHandler)
	router.GET("/list", s.ListHandler)
	router.POST("/sessions/new", s.NewSessionHandler)
	router.GET("/sessions/list", s.ListSessionHandler)
	router.GET("/session/:id", s.SessionHandler)
	router.GET("/session_choice/get", s.SessionChoiceGetHandler)
	router.POST("/session_choice/keep", s.SessionChoiceKeepHandler)
	router.POST("/session_choice/remove", s.SessionChoiceRemoveHandler)
	router.GET("/session_choice/list", s.SessionChoiceListHandler)
	router.NotFound = http.HandlerFunc(s.NotFoundHandler)
}

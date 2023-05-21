package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/superchausette/le-prenom/leprenom"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	host := flag.String("host", "", "host to bind to")
	port := flag.Int("port", 9999, "port to listen on")
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

	server := leprenom.NewServer(db)
	router := httprouter.New()
	server.SetupRoutes(router)

	// Start the web server
	addr := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Println("Server is listening on ", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

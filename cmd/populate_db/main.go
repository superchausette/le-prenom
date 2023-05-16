package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/superchausette/le-prenom/leprenom"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	csvFile := flag.String("csv", "", "CSV file to use")
	dbName := flag.String("dbname", "", "Database to create or update")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if *csvFile == "" {
		flag.Usage()
		fmt.Println("Missing csv file argument")
		return
	}
	if *dbName == "" {
		flag.Usage()
		fmt.Println("Missing database name argument")
		return
	}

	fmt.Println("Database name:", *dbName)
	// Your database connection logic here
	data, err := os.Open(*csvFile)
	if err != nil {
		panic("Unable to open data/nat2021.csv")
	}

	// Retrieving data to insert
	fmt.Println("Reading CSV file ", *csvFile)
	csvContent, err := leprenom.Import(data)
	firstNameToInsert := []leprenom.FirstName{}
	for _, content := range csvContent {
		firstNameToInsert = append(firstNameToInsert, leprenom.FirstName{Name: content.FirstName, Gender: content.Gender})
	}

	fmt.Println("Opening sqlite database ", *dbName)
	db, err := gorm.Open(sqlite.Open(*dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&leprenom.FirstName{})

	fmt.Println("Creating ", len(firstNameToInsert), " transactions ")
	// Insert first name in the db using a transaction
	transaction := db.Begin()
	for _, toInsert := range firstNameToInsert {
		result := transaction.Create(&toInsert)
		if result != nil && result.Error != nil {
			fmt.Println("Error for ", toInsert, " -> ", result.Error)
		}

	}
	fmt.Println("Commit transaction to DB")
	transaction.Commit()

	// Print the number of entry in the db
	firstNameCnt, err1 := leprenom.CountFirstName(db)
	if err1 != nil {
		panic(err)
	}
	boyFirstNameCnt, err2 := leprenom.CountBoyFirstName(db)
	if err2 != nil {
		panic(err)
	}
	girlFirstNameCnt, err3 := leprenom.CountGirlFirstName(db)
	if err3 != nil {
		panic(err)
	}
	fmt.Println(firstNameCnt, "first name in database")
	fmt.Println(boyFirstNameCnt, "girl first name in database")
	fmt.Println(girlFirstNameCnt, "boy first name in database")
}

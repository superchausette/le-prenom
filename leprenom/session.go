package leprenom

import (
	"fmt"

	"gorm.io/gorm"
)

func SessionTypeFromString(value string) uint {
	switch value {
	case "all":
		return AllSession
	case "boy":
		return BoySession
	case "girl":
		return GirlSession
	case "unisex":
		return UnisexSession
	}
	panic(fmt.Sprintf("Unknown session type '%s'", value))
}

func SessionTypeToString(value uint) string {
	switch value {
	case AllSession:
		return "all"
	case BoySession:
		return "boy"
	case GirlSession:
		return "girl"
	case UnisexSession:
		return "unisex"
	}
	panic(fmt.Sprintf("Unknown session type '%d'", value))
}

func CreateSession(sessionName, sessionType string, db *gorm.DB) string {
	fmt.Println("Creating session: ", sessionName)
	result := db.Create(&Session{Name: sessionName, FirstNameType: SessionTypeFromString(sessionType)})
	if result != nil && result.Error != nil {
		return fmt.Sprintf("Unable to create session '%s' '%s'", sessionName, sessionType)
	}
	return fmt.Sprintf("Successfully created session '%s' '%s'", sessionName, sessionType)
}

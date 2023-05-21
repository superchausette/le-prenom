package leprenom

import (
	"fmt"
	"html/template"
	"net/http"

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

func TotalFirstNameAvailableForSession(sessionType uint, db *gorm.DB) int64 {
	switch sessionType {
	case AllSession:
		count, err := CountFirstName(db)
		if err != nil {
			fmt.Println("Unable to count first name: ", err)
			return 0
		}
		return count
	case BoySession:
		count, err := CountBoyFirstName(db)
		if err != nil {
			fmt.Println("Unable to count first name: ", err)
			return 0
		}
		return count
	case GirlSession:
		count, err := CountGirlFirstName(db)
		if err != nil {
			fmt.Println("Unable to count first name: ", err)
			return 0
		}
		return count
	case UnisexSession:
		count, err := CountUnisexName(db)
		if err != nil {
			fmt.Println("Unable to count first name: ", err)
			return 0
		}
		return count
	}
	panic(fmt.Sprintf("Unknown session type '%d'", sessionType))
}

func RemainingFirstNameAvailableForSession(sessionId, sessionType uint, db *gorm.DB) int64 {
	var count int64
	switch sessionType {
	case AllSession:
		err := db.Table("first_names").
			Select("COUNT(DISTINCT first_names.name)"). // TODO issue with unisex names.
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				" AND session_contents.session_id = ?", sessionId).
			Where("session_contents.first_name_id IS NULL").
			Count(&count).Error
		if err != nil {
			fmt.Println("Unable to count remaining first name for session ", sessionId, ": ", err)
			return 0
		}
		return count
	case BoySession:
		err := db.Table("first_names").
			Select("COUNT(DISTINCT first_names.name)").
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				"AND session_contents.session_id = ?", sessionId).
			Where("first_names.gender = ? and session_contents.first_name_id IS NULL", BoyFirstName).
			Count(&count).Error
		if err != nil {
			fmt.Println("Unable to count remaining first name for session ", sessionId, ": ", err)
			return 0
		}
		return count
	case GirlSession:
		err := db.Table("first_names").
			Select("COUNT(DISTINCT first_names.name)").
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				"AND session_contents.session_id = ?", sessionId).
			Where("first_names.gender = ? and session_contents.first_name_id IS NULL", GirlFirstName).
			Count(&count).Error
		if err != nil {
			fmt.Println("Unable to count remaining first name for session ", sessionId, ": ", err)
			return 0
		}
		return count
		/*
			case UnisexSession:
				count, err := CountUnisexName(db)
				if err != nil {
					fmt.Println("Unable to count first name: ", err)
					return 0
				}
				return count
			}
		*/
	}
	panic(fmt.Sprintf("Unknown session type '%d'", sessionType))
}

func GetRandomRemainingFirstNameAvailableForSession(sessionId, sessionType uint, db *gorm.DB) FirstName {
	var name FirstName
	switch sessionType {
	case AllSession:
		err := db.Table("first_names").
			Select("first_names.*"). //TODO issue with unisex names.
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				" AND session_contents.session_id = ?", sessionId).
			Where("session_contents.first_name_id IS NULL").
			Order("RANDOM()").
			First(&name).Error
		if err != nil {
			fmt.Println("Unable to get remaining first name for session ", sessionId, ": ", err)
			return name
		}
		return name
	case BoySession:
		err := db.Table("first_names").
			Select("first_names.*").
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				"AND session_contents.session_id = ?", sessionId).
			Where("first_names.gender = ? and session_contents.first_name_id IS NULL", BoyFirstName).
			Order("RANDOM()").
			First(&name).Error
		if err != nil {
			fmt.Println("Unable to get remaining first name for session ", sessionId, ": ", err)
			return name
		}
		return name
	case GirlSession:
		err := db.Table("first_names").
			Select("first_names.*").
			Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id",
				"AND session_contents.session_id = ?", sessionId).
			Where("first_names.gender = ? and session_contents.first_name_id IS NULL", GirlFirstName).
			Order("RANDOM()").
			First(&name).Error
		if err != nil {
			fmt.Println("Unable to get remaining first name for session ", sessionId, ": ", err)
			return name
		}
		return name
		/*
			case UnisexSession:
				count, err := CountUnisexName(db)
				if err != nil {
					fmt.Println("Unable to count first name: ", err)
					return 0
				}
				return count
			}
		*/
	}
	panic(fmt.Sprintf("Unknown session type '%d'", sessionType))
}

func SessionFirstNameKept(sessionID uint, db *gorm.DB) []string {
	var firstNames []FirstName
	err := db.
		Table("first_names").
		Select("first_names.name").
		Joins("LEFT JOIN session_contents ON first_names.id = session_contents.first_name_id").
		Where("session_contents.session_id = ? AND session_contents.keep = ?", sessionID, true).
		Find(&firstNames).
		Error
	if err != nil {
		fmt.Println("Unable to get kept name for ", sessionID, ": ", err)
		return []string{}
	}
	ret := make([]string, len(firstNames))
	for idx, entry := range firstNames {
		ret[idx] = entry.Name
	}
	return ret
}

func WriteSessionFirstNameEntry(sessionID, tableIndex uint,
	CallerName string,
	w http.ResponseWriter,
	db *gorm.DB,
	template *template.Template) {
	var session Session
	sessionErr := db.First(&session, sessionID).Error
	if sessionErr != nil {
		fmt.Println(CallerName, ": Unable to query session ID ", sessionID, " - ", sessionErr)
		http.Error(w, fmt.Sprint(CallerName, ": Unable to query session ID", sessionID, sessionErr), http.StatusBadRequest)
		return
	}
	randomFirstName := GetRandomRemainingFirstNameAvailableForSession(session.ID, session.FirstNameType, db)
	tmplParam := struct {
		SessionID   uint
		FirstNameID uint
		FirstName   string
		Index       uint
	}{session.ID,
		randomFirstName.ID,
		randomFirstName.Name,
		tableIndex}
	tmplErr := template.ExecuteTemplate(w, "session_first_name_table_entry.html", tmplParam)

	if tmplErr != nil {
		fmt.Println(CallerName, ": unable to execute template", tmplErr)
	}
}

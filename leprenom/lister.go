package leprenom

import (
	"fmt"

	"gorm.io/gorm"
)

func ListAllFirstName(db *gorm.DB) []string {
	var first_names []string
	result := db.Model(&FirstName{}).Select("name").Find(&first_names)
	if result.Error != nil {
		fmt.Println("Unable to get all first name count", result.Error)

	}
	return first_names
}

func ListAllBoyFirstName(db *gorm.DB) []string {
	var first_names []string
	result := db.Model(&FirstName{}).Select("name").Where("Gender = ?", "1").Find(&first_names)
	if result.Error != nil {
		fmt.Println("Unable to get all boy first name", result.Error)

	}
	return first_names
}

func ListAllGirlFirstName(db *gorm.DB) []string {
	var first_names []string
	result := db.Model(&FirstName{}).Select("name").Where("Gender = ?", "2").Find(&first_names)
	if result.Error != nil {
		fmt.Println("Unable to get all girl first name", result.Error)

	}
	return first_names
}

func ListAllUnisexFirstName(db *gorm.DB) []string {
	var first_names []string
	result := db.Model(&FirstName{}).Select("name").Group("name").Having("COUNT(DISTINCT gender) > ?", 1).Find(&first_names)
	if result.Error != nil {
		fmt.Println("Unable to get all unisex first name", result.Error)

	}
	return first_names
}

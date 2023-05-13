package leprenom

import "gorm.io/gorm"

type FirstName struct {
	gorm.Model
	Name   string `gorm:"uniqueIndex:idx_name_gender;size:128"`
	Gender int    `gorm:"uniqueIndex:idx_name_gender"`
}

type Session struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex;size:64"`
	content []SessionContent
}

type SessionContent struct {
	gorm.Model
	FirstNameID uint
	Status      string
}

package leprenom

import (
	"time"
)

const (
	BoyFirstName  int = 1
	GirlFirstName     = 2
)

const (
	AllSession    uint = 1
	BoySession         = 2
	GirlSession        = 3
	UnisexSession      = 4
)

type FirstName struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"uniqueIndex:idx_name_gender;type:varchar(64)"`
	Gender int    `gorm:"uniqueIndex:idx_name_gender"`
}

type Session struct {
	ID              uint   `gorm:"primaryKey;autoIncrement"`
	Name            string `gorm:"type:varchar(64);unique"`
	FirstNameType   uint   `gorm:"default:1"`
	SessionContents []SessionContent
}

type SessionContent struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	SessionID        uint `gorm:"uniqueIndex:idx_name_session"`
	FirstNameID      uint `gorm:"uniqueIndex:idx_name_session"`
	Keep             bool
	LastModifiedTime time.Time `gorm:"autoUpdateTime:milli"`
}

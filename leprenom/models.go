package leprenom

type FirstName struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Name   string `gorm:"uniqueIndex:idx_name_gender;type:varchar(64)"`
	Gender int    `gorm:"uniqueIndex:idx_name_gender"`
}

type Session struct {
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	Name    string `gorm:"type:varchar(64);unique"`
	content []SessionContent
}

type SessionContent struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	FirstNameID uint
	StatusID    uint
}

type SessionNameStatus struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Status string `gorm:"type:varchar(32);unique"`
}

package leprenom

import "gorm.io/gorm"

func CountFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Count(&count)
	return count, result.Error
}

func CountBoyFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Where("Gender = ?", BoyFirstName).Count(&count)
	return count, result.Error
}

func CountGirlFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Where("Gender = ?", GirlFirstName).Count(&count)
	return count, result.Error
}

func CountUnisexName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Select("name").Group("name").Having("COUNT(DISTINCT gender) > ?", 1).Count(&count)
	return count, result.Error
}

package leprenom

import "gorm.io/gorm"

func CountFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Count(&count)
	return count, result.Error
}

func CountBoyFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Where("Gender = ?", "1").Count(&count)
	return count, result.Error
}

func CountGirlFirstName(db *gorm.DB) (int64, error) {
	var count int64
	result := db.Model(&FirstName{}).Where("Gender = ?", "1").Count(&count)
	return count, result.Error
}

package leprenom

import (
	"fmt"

	"gorm.io/gorm"
)

type FirstNameStats struct {
	Count       int64
	BoyCount    int64
	GirlCount   int64
	UnisexCount int64
}

func NewFirstNameStats(db *gorm.DB) *FirstNameStats {
	// Print the number of entry in the db
	var stats FirstNameStats
	{
		cnt, err := CountFirstName(db)
		if err != nil {
			fmt.Println("Unable to get first name count", err)
		}
		stats.Count = cnt
	}
	{
		cnt, err := CountBoyFirstName(db)
		if err != nil {
			fmt.Println("Unable to get boy first name count", err)
		}
		stats.BoyCount = cnt
	}
	{
		cnt, err := CountGirlFirstName(db)
		if err != nil {
			fmt.Println("Unable to get boy first name count", err)
		}
		stats.GirlCount = cnt
	}
	{
		cnt, err := CountUnisexName(db)
		if err != nil {
			fmt.Println("Unable to get boy first name count", err)
		}
		stats.UnisexCount = cnt
	}

	return &stats
}

package scripts

import (
	"fmt"

	"gorm.io/gorm"
)

func FindUserEmail(db *gorm.DB, formEmail string) bool {

	var flag bool
	var checksID uint64

	db.Table("users").Select("id").Where("email = ?", formEmail).First(&checksID)
	if checksID == 0 {
		flag = false
	} else {
		flag = true
	}
	fmt.Println(checksID)

	return flag
}

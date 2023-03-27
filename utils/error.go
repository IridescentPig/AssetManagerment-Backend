package utils

import (
	"log"

	"gorm.io/gorm"
)

func DB_error(result *gorm.DB) error {
	err := result.Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func Service_error(code int, info string) (int, string) {
	log.Println(info)
	return code, info
}

package utils

import (
	"log"

	"gorm.io/gorm"
)

func DBError(result *gorm.DB) error {
	err := result.Error
	if err != nil {
		log.Println(err)
	}
	return err
}

func ServiceError(code int, info string) (int, string) {
	log.Println(info)
	return code, info
}

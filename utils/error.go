package utils

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

func DBError(result *gorm.DB) error {
	err := result.Error
	if err != nil {
		log.Println(err)
		return errors.New("database error")
	}
	return nil
}

func ServiceError(code int, info string) (int, string) {
	log.Println(info)
	return code, info
}

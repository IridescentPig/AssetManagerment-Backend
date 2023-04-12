package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const salt = "se-asset-BinaryAbstract-2023"

func CreateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str + salt))
	return hex.EncodeToString(h.Sum(nil))
}

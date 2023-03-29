package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const saltPassword = "se-asset-BinaryAbstract-2023"

func CreateMD5(str string) string {
	h := md5.New()
	h.Write([]byte(str + saltPassword))
	return hex.EncodeToString(h.Sum(nil))
}

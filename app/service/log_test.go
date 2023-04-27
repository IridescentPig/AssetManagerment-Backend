package service

import "testing"

func TestLog(t *testing.T) {
	LogService.GetLoginLog(1)
	LogService.GetLoginLog(9)
	LogService.GetDataLog(1)
	LogService.GetDataLog(9)
}

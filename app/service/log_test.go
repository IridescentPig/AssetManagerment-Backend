package service

import "testing"

func TestLog(t *testing.T) {
	LogService.GetLoginLog(1, 50, 0)
	LogService.GetLoginLog(9, 50, 0)
	LogService.GetDataLog(1, 50, 0)
	LogService.GetDataLog(9, 50, 0)
}

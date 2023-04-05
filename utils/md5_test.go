package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	result := CreateMD5("123456")
	assert.Equal(t, "dc171ced6d63ab8021656851788801c7", result, "md5 error")
}

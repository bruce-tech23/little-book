package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("1234#test") // 注意：bcrypt 加密字符串长度不超过72字节
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	fmt.Println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("1234#test"))
	assert.NoError(t, err)
}

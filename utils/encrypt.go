package utils

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
)

// 用户密码的加密与验证
func EncryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 生成验证码
var verifyCodeTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GenerateVerifyCode(length int) (string, error) {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if err != nil || n != length {
		return "", errors.New("生成验证码出错")
	}
	for i := 0; i < len(b); i++ {
		b[i] = verifyCodeTable[int(b[i])%len(verifyCodeTable)]
	}
	return string(b), nil
}

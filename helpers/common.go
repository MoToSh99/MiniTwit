package helpers

import (
    "time"
    "golang.org/x/crypto/bcrypt"
)

func IsEmpty(data string) bool {
    if len(data) <= 0 {
        return true
    } else {
        return false
    }
}

func GetCurrentTime() string{
	dt := time.Now()
	return (dt.Format("15:04:05 02-01-2006"))
}

func HashPassword(password string) (string) {
    bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
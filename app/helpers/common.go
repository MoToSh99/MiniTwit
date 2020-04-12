package helpers

import (
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	
	"golang.org/x/crypto/bcrypt"
	"os"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

func IsEmpty(data string) bool {
	if len(data) <= 0 {
		return true
	} else {
		return false
	}
}

func GetCurrentTime() string {
	dt := time.Now()
	return (dt.Format("2006-01-02 15:04:05"))
}

func GetConnString() string {

	var connString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PASSWORD"))
	return connString
}

func GetDB() *gorm.DB {
	return db
}

func InitDB() *gorm.DB {

	var connString = GetConnString()

	db, err = gorm.Open("postgres", connString)

	db.DB().SetMaxIdleConns(25)
	db.DB().SetMaxOpenConns(25)
	db.DB().SetConnMaxLifetime(5*time.Minute)



	if err != nil {
		panic("failed to connect database")
	}

	return db

}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

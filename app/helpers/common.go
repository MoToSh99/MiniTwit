package helpers

import (
	"fmt"
	"time"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"os"
	"log"
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

	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Error loading .env file")
	}

	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		os.Getenv("SERVER_ADDR"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("SERVER_PORT"), os.Getenv("DATABASE_NAME_PUB"))
	return connString
}

func GetDB() *gorm.DB {
	return db
}

func InitDB() *gorm.DB {

	var connString = GetConnString()

	db, err = gorm.Open("mssql", connString)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	//db.DB().SetMaxIdleConns(0)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(500)

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

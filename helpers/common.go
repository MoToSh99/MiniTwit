package helpers

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"golang.org/x/crypto/bcrypt"
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

	var server = "minitwitserver.database.windows.net"
	var port = 1433
	var user = "Minitwit"
	var password = "ITU2020!"
	var database = "minitwitdb"

	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	return connString
}

func GetDB() *gorm.DB {

	var server = "minitwitserver.database.windows.net"
	var port = 1433
	var user = "Minitwit"
	var password = "ITU2020!"
	var database = "minitwitdb"

	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	db, err := gorm.Open("sqlite3", "minitwit.db")

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

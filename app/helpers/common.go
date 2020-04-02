package helpers

import (
	"fmt"
	"time"

	c "../config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/spf13/viper"
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
	// Set the file name of the configurations file
	viper.SetConfigName("config")

	// Set the path to look for the configurations file
	viper.AddConfigPath("../")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var configuration c.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		viper.GetString("server.addr"), viper.GetString("database.user"), viper.GetString("database.password"), viper.GetInt("server.port"), viper.GetString("database.name_pub"))
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

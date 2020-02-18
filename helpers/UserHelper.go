package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	cookies "../cookies"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	structs "../structs"
)

var databasepath = "/tmp/minitwit.db"

type User struct {
	user_id   int
	username  string
	email     string
	pw_hash   string
	image_url string
}

type Post struct {
	Username      string
	PostMessageid int
	AuthorId      int
	Text          string
	Date          string
	Flag          int
	Image         string
}

func GetUserID(username string) int {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	
	user := structs.User{}

	db.Where("username = ?", username).First(&user)

	return user.User_id

}

func CheckUsernameExists(username string) bool {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	type Output struct {
		Id int
	}
	var output []Output
	db.Table("users").Select("user_id").Where("username = ?", username).Scan(&output)
	var rtn bool
	// Catch errors
	if len(output) < 1 {
		rtn = false
	} else {
		rtn = true
	}
	return rtn

}

func GetGravatarHash(g_email string) string {
	g_email = strings.TrimSpace(g_email)
	g_email = strings.ToLower(g_email)
	h := md5.New()
	io.WriteString(h, g_email)
	finalBytes := h.Sum(nil)
	finalString := hex.EncodeToString(finalBytes)
	return finalString
}

func GetAllPosts() []Post {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var post []Post
	db.Table("users").Select("users.username").Joins("join messages on users.user_id = messages.author_id").Order("messages.pub_date desc").Scan(&post)

	return post

}

func GetUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("cookie"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookies.CookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func ValidUser(username string, psw string) bool {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	type Output struct {
		Id int
	}
	var output []Output
	db.Table("users").Select("user_id").Where("username = ? AND pw_hash = ?", username, psw).Scan(&output)
	var rtn bool
	// Catch errors
	if len(output) < 1 {
		rtn = false
	} else {
		rtn = true
	}
	return rtn
}

func UserIsValid(uName, pwd string) bool {
	_isValid := false

	if ValidUser(uName, pwd) {
		_isValid = true
	} else {
		_isValid = false
	}

	return _isValid
}

func GetUserPosts(username string) []Post {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var post []Post

	db.Table("users").Select("users.username, m.*").Joins("join messages m on users.user_id = m.author_id").Where("m.flagged = ? and m.author_id = u.user_id and (u.user_id = ? or u.user_id in (select whom_id from followers where who_id = ?))", 0, GetUserID(username), GetUserID(username)).Scan(&post)

	return post

}

func GetUsernameFromID(id int) string {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	user := structs.User{}

	db.Where("user_id = ?", id).First(&user)

	return user.Username

}

func PostsAmount(posts []Post) bool {

	if len(posts) > 0 {
		return true
	} else {
		return false
	}
}

func CheckIfFollowed(who string, whom string) bool {
	db, err := gorm.Open("sqlite3", databasepath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	type Output struct {
		Who_id  int
		Whom_id int
	}
	var output []Output
	db.Table("followers").Select("*").Where("followers.who_id = ? AND followers.whom_id = ?", who, whom).Scan(&output)
	var rtn bool
	// Catch errors
	if len(output) < 1 {
		rtn = false
	} else {
		rtn = true
	}
	return rtn
}

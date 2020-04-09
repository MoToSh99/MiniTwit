package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	cookies "../cookies"
	structs "../structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB
var err error

type Post struct {
	Username   string
	Message_id int
	Author_id  int
	Text       string
	Pub_date   string
	Flagged    int
	Image_url  string
}

func GetUserID(username string) int {
	db := GetDB()

	user := structs.User{}

	db.Where("username = ?", username).First(&user)

	return user.User_id

}

func CheckUsernameExists(username string) bool {
	db := GetDB()

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
	db := GetDB()

	var postSlice []Post
	db.Table("messages").Offset(10).Limit(10).Order("messages.pub_date desc").Select("users.username, messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged, users.image_url").Joins("join users on users.user_id = messages.author_id").Where("messages.flagged = 0").Scan(&postSlice)

	return postSlice
}

func GetMoreposts(numberOfPosts int) []Post {
	db := GetDB()

	var posts []structs.Message
	
	db.Where("flagged = ?", 0).Limit(numberOfPosts).Order("pub_date desc").Find(&ms)
	
	var postSlice []Post

	for _, m := range posts {
		if err != nil {
			continue
		}
		post := Post{
			Username:      GetUsernameFromID(m.Author_id),    
			Message_id:	   m.Message_id,
			Author_id:	   m.Author_id,
			Text:          m.Text,
			Pub_date:	   m.Pub_date,
			Flagged:       m.Flagged,
			Image_url:     GetImageFromID(m.Author_id),
		}
		postSlice = append(postSlice, post)
	}
	return postSlice
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
	db := GetDB()

	user := structs.User{}

	db.Where("username = ?", username).First(&user)

	if CheckPasswordHash(psw, user.Pw_hash) {
		return true
	} else {
		return false
	}

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
	db := GetDB()

	var postSlice []Post

	db.Table("messages").Order("messages.pub_date desc").Select("users.username, messages.message_id, messages.author_id, messages.text, messages.pub_date, messages.flagged, users.image_url").Joins("join users on users.user_id = messages.author_id").Where("messages.flagged = 0 AND users.username = ?", username).Scan(&postSlice)
	return postSlice

}

func GetUsernameFromID(id int) string {
	db := GetDB()

	user := structs.User{}

	db.Where("user_id = ?", id).First(&user)

	return user.Username

}

func GetImageFromID(id int) string {
	db := GetDB()

	user := structs.User{}

	db.Where("user_id = ?", id).First(&user)

	return user.Image_url

}

func PostsAmount(posts []Post) bool {

	if len(posts) > 0 {
		return true
	} else {
		return false
	}
}

func CheckIfFollowed(who string, whom string) bool {
	db := GetDB()
	output := []structs.Follower{}
	db.Where("who_id = ? AND whom_id = ?", GetUserID(whom), GetUserID(who)).Scan(&output)

	var rtn bool
	if len(output) < 1 {
		rtn = false
	} else {
		rtn = true
	}

	return rtn
}

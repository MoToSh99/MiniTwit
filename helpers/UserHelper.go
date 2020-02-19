package helpers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	cookies "../cookies"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	structs "../structs"
	"database/sql"
)

var db *sql.DB

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
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	
	user := structs.User{}

	db.Where("username = ?", username).First(&user)

	return user.User_id

}

func CheckUsernameExists(username string) bool {
	db, err := gorm.Open("mssql", GetConnString())
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
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	messages := []structs.Message{}

	db.Order("pub_date desc").Where("flagged = ?",0).Find(&messages)

	var postSlice []Post
	for _, m := range messages {
		post := Post{Username: GetUsernameFromID(m.Author_id), PostMessageid: m.Message_id, Text: m.Text, Date: m.Pub_date, Image: GetImageFromID(m.Author_id)  }
		postSlice = append(postSlice, post)
	}

	return postSlice;
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
	db, err := gorm.Open("mssql", GetConnString())
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
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	messages := []structs.Message{}
	
	db.Order("pub_date desc").Where("flagged = ? AND author_id = ? ",0, GetUserID(username)).Or("author_id in (select whom_id from followers where who_id = ?)", GetUserID(username)).Find(&messages)

	var postSlice []Post
	for _, m := range messages {
		post := Post{Username: GetUsernameFromID(m.Author_id), PostMessageid: m.Message_id, Text: m.Text, Date: m.Pub_date, Image: GetImageFromID(m.Author_id)  }
		postSlice = append(postSlice, post)
	}

	return postSlice;

}



func GetUsernameFromID(id int) string {
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	user := structs.User{}

	db.Where("user_id = ?", id).First(&user)

	return user.Username

}

func GetImageFromID(id int) string {
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

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
	db, err := gorm.Open("mssql", GetConnString())
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()	
	output := []structs.Follower{}
	db.Where("who_id = ? AND whom_id = ?", GetUserID(whom), GetUserID(who)).Find(&output)

	var rtn bool
	if len(output) < 1 {
		rtn = false
	} else {
		rtn = true
	}

	return rtn
}

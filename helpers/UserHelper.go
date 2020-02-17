package helpers

import (
	"strings"
	"fmt"
	"io"
	"encoding/hex"
	"crypto/md5"
	"net/http"
	cookies "../cookies"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)


var databasepath = "/tmp/minitwit.db"
var database, _ = gorm.Open("sqlite3", databasepath)

type id struct{
	id int
}

type user struct  {
	user_id int
	username string
	email string
	pw_hash string
	image_url string
  }

  type Post struct {
	Username string
	PostMessageid int
	AuthorId int 
	Text string
	Date string
	Flag int
	Image string
}

func GetUserID(username string) int{
	var output id
	database.Raw("SELECT user_id FROM users WHERE username = ?", username).Scan(&output)
	return output.id

}

func CheckUsernameExists(username string) bool {
	var output bool
	
	var count int
	// Prepare your query
	database.Raw("SELECT user_id FROM users WHERE username = ?", username).Count(&count)

	// Catch errors
	switch {
	case count == 0:
			output = false
	default:
			output = true
	}

	return output

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




func GetAllPosts()[]Post{
	var post Post
	//sqlStatement := `SELECT u.username, m.message_id, m.text, m.pub_date, u.image_url FROM message m join user u ON m.author_id = u.user_id order by m.pub_date desc`
	rows, err := database.Table("users").Select("users.username").Joins("join messages on users.user_id = messages.author_id").Order("messages.pub_date desc").Rows()
	//rows, err := database.Query(sqlStatement)
	if err != nil {
		panic(err)
	   }
	
	defer rows.Close()


	var postSlice []Post
	for rows.Next(){
		rows.Scan(&post.Username, &post.PostMessageid, &post.Text, &post.Date, &post.Image)
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
	var output bool


	var database, _ = gorm.Open("sqlite3", databasepath)

	var count int
	// Prepare your query
	database.Raw("SELECT user_id FROM users WHERE username = ? AND pw_hash = ?", username, psw).Count(&count)

	// Catch errors
	switch {
	case count == 0:
			output = false
	default:
			output = true
	}

	return output

}

func UserIsValid(uName, pwd string) bool {
    _isValid :=  false
 
    if ValidUser(uName,pwd) {
        _isValid = true
    } else {
        _isValid = false
    }
 
    return _isValid
}

func GetUserPosts(username string)[]Post{
	var post Post


	var database, _ = gorm.Open("sqlite3", databasepath)
	//rows, _ := database.Raw("select m.*, u.image_url  from message m JOIN user u on m.author_id = u.user_id where m.flagged = 0 and m.author_id = u.user_id and (u.user_id = ? or u.user_id in (select whom_id from follower where who_id = ?)) order by m.pub_date desc",GetUserID(username),GetUserID(username)).Rows()


	rows, err := database.Table("users").Select("users.username, m.*").Joins("join messages m on users.user_id = m.author_id").Where("m.flagged = ? and m.author_id = u.user_id and (u.user_id = ? or u.user_id in (select whom_id from followers where who_id = ?))",0,GetUserID(username),GetUserID(username)).Rows()
	
	if (err != nil){
		panic(err)
	}
	
	var postSlice []Post
	for rows.Next() {
		rows.Scan(&post.PostMessageid, &post.AuthorId, &post.Text, &post.Date, &post.Flag, &post.Image)
		post.Username = GetUsernameFromID(post.AuthorId)
		fmt.Sprintf(post.Username)
		postSlice = append(postSlice, post)
	}

	

	database.Close()
	
	return postSlice;

}

func GetUsernameFromID(id int) string{
	var output user

	var database, _ = gorm.Open("sqlite3", databasepath)
	database.Raw("SELECT username FROM users WHERE user_id = ?", id).Find(&output)
	
	database.Close()

	return output.username

}



func PostsAmount(posts []Post) bool{

	if (len(posts) > 0){
		return true
	} else {
		return false
	}
}


func CheckIfFollowed(who string, whom string) bool{
	var output bool

	// Prepare your query
	var database, _ = gorm.Open("sqlite3", databasepath)
	var count int
	database.Raw("select * from followers where followers.who_id = ? and followers.whom_id = ?",who,whom).Count(&count)
	
	switch {
	case count == 0:
			output = false
	default:
			output = true
	}
	database.Close()
	return output
}

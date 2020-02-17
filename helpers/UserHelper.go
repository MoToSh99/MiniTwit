package helpers

import (
	"database/sql"
	"fmt"
	"strings"
	"io"
	"encoding/hex"
	"crypto/md5"
	"net/http"
	cookies "../cookies"
)


var databasepath = "/tmp/minitwit.db"
var database, _ = sql.Open("sqlite3", databasepath)


func GetUserID(username string) int{
	var output int
	// Prepare your query
	query, err := database.Prepare("SELECT user_id FROM user WHERE username = ?")

	if err != nil {
		fmt.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(username).Scan(&output)

	return output

}

func CheckUsernameExists(username string) bool {
	var output bool

	// Prepare your query
	query, err := database.Prepare("SELECT user_id FROM user WHERE username= ?")

	if err != nil {
		fmt.Printf("%s", err)
	}

	defer query.Close()

	err = query.QueryRow(username).Scan(&output)

	// Catch errors
	switch {
	case err == sql.ErrNoRows:
			output = false
	case err != nil:
			output = true
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

type Post struct {
	Username string
	PostMessageid int
	AuthorId int 
	Text string
	Date string
	Flag int
	Image string
}


func GetAllPosts()[]Post{
	var post Post
	sqlStatement := `SELECT u.username, m.message_id, m.text, m.pub_date, u.image_url FROM message m join user u ON m.author_id = u.user_id order by m.pub_date desc`
	rows, err := database.Query(sqlStatement)
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


	var database, _ = sql.Open("sqlite3", databasepath)
	// Prepare your query
	query, err := database.Prepare("SELECT user_id FROM user WHERE username = ? AND pw_hash = ?")

	
	if err != nil {
		fmt.Printf("%s", err)
	}

	defer query.Close()

	err = query.QueryRow(username, psw).Scan(&output)

	// Catch errors
	switch {
	case err == sql.ErrNoRows:
			output = false
	case err != nil:
			output = true
	default:
			output = true
	}

	database.Close()

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


	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("select m.*, u.image_url  from message m JOIN user u on m.author_id = u.user_id where m.flagged = 0 and m.author_id = u.user_id and (u.user_id = ? or	u.user_id in (select whom_id from follower where who_id = ?)) order by m.pub_date desc")

	if err != nil {
		fmt.Printf("%s", err)
	}


	rows, err := query.Query(GetUserID(username),GetUserID(username))
	defer rows.Close()
	
	var postSlice []Post
	for rows.Next(){
		rows.Scan(&post.PostMessageid, &post.AuthorId, &post.Text, &post.Date, &post.Flag, &post.Image)
		post.Username = GetUsernameFromID(post.AuthorId)
		postSlice = append(postSlice, post)
	}

	database.Close()
	
	return postSlice;

}

func GetUsernameFromID(id int) string{
	var output string

	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("SELECT username FROM user WHERE user_id = ?")
	

	if err != nil {
		fmt.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(id).Scan(&output)

	database.Close()

	return output

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
	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("select * from follower where follower.who_id = ? and follower.whom_id = ?")
	

	if err != nil {
		fmt.Printf("%s", err)
	}

	defer query.Close()

	err = query.QueryRow(GetUserID(whom), GetUserID(who)).Scan(&output)

	// Catch errors
	switch {
	case err == sql.ErrNoRows:
			output = false
	case err != nil:
			output = true
	default:
			output = true
	}
	database.Close()

	return output
}

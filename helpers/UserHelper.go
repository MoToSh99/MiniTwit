package helpers

import (
	"database/sql"
	"fmt"
	"strings"
	"io"
	"encoding/hex"
	"crypto/md5"
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


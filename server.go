package main

import (
	"github.com/gorilla/mux"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/gorilla/securecookie"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"strings"
	"io"
	"encoding/hex"
	"crypto/md5"
	api "./api"
)

var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
	
var templates map[string]*template.Template

func init() {
	loadTemplates()
}

var databasepath = "/tmp/minitwit.db"

func main() {

	var database, _ = sql.Open("sqlite3", databasepath)
	statement, _ := database.Prepare("create table if not exists user (user_id integer primary key autoincrement,username string not null,email string not null,pw_hash string not null,image_url string);")
	statement.Exec()
	statement2, _ := database.Prepare("create table if not exists follower ( who_id integer, whom_id integer);")
	statement2.Exec()
	statement3, _ := database.Prepare("create table if not exists message (message_id integer primary key autoincrement,author_id integer not null,text string not null,pub_date integer,flagged integer);")
	statement3.Exec()
	database.Close()

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	router.HandleFunc("/", PublicTimelineRoute).Methods("GET")
	router.HandleFunc("/about", AboutRoute).Methods("GET")
	router.HandleFunc("/contact", ContactRoute).Methods("GET")

	router.HandleFunc("/signin", SigninRoute).Methods("GET")
	router.HandleFunc("/signin", LoginHandler).Methods("POST")


	router.HandleFunc("/register", SignupRoute).Methods("GET")
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	
	router.HandleFunc("/personaltimeline", PersonalTimelineRoute).Methods("GET")
	router.HandleFunc("/personaltimeline", PersonalTimelineHandler).Methods("POST")
	

	router.HandleFunc("/signout", LogoutHandler)

	router.HandleFunc("/publictimeline", PublicTimelineRoute).Methods("GET")
	router.HandleFunc("/publictimeline", PublicTimelineHandler).Methods("POST")
	

	router.HandleFunc("/{username}", UserpageRoute).Methods("GET")
	router.HandleFunc("/{username}", UserpageHandler).Methods("POST")

	router.HandleFunc("/{username}/follow", UserFollowHandler)
	router.HandleFunc("/{username}/unfollow", UserUnfollowHandler)

	

	apiRoute := mux.NewRouter()
	apiRoute.HandleFunc("/test", api.Test)
	apiRoute.HandleFunc("/latest", api.Get_latest)
	apiRoute.HandleFunc("/register", api.Register).Methods("POST")
	apiRoute.HandleFunc("/msgs", api.Messages)
	apiRoute.HandleFunc("/msgs/{username}", api.Messages_per_user)
	apiRoute.HandleFunc("/fllws/{username}", api.Follow)



	port := 4999
	log.Printf("Server starting on port %v\n", port)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))}()

	apiport := 5001
	log.Printf("Api Server starting on port %v\n", apiport)
    go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", apiport), apiRoute))}()
	
    select {}
}



func UserIsValid(uName, pwd string) bool {
    _isValid :=  false
 
    if validUser(uName,pwd) {
        _isValid = true
    } else {
        _isValid = false
    }
 
    return _isValid
}

func IsEmpty(data string) bool {
    if len(data) <= 0 {
        return true
    } else {
        return false
    }
}

func IndexRoute(res http.ResponseWriter, req *http.Request) {

	if err := templates["publictimeline"].Execute(res, map[string]interface{}{
        "loggedin": !IsEmpty(GetUserName(req)),
    }); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AboutRoute(res http.ResponseWriter, req *http.Request) {

	if err := templates["about"].Execute(res, map[string]interface{}{
        "loggedin": !IsEmpty(GetUserName(req)),
    }); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func UserFollowHandler(res http.ResponseWriter, req *http.Request){
	if (IsEmpty(GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	var database, _ = sql.Open("sqlite3", databasepath)
	statement, _ := database.Prepare("insert into follower (who_id, whom_id) values (?, ?)")
	statement.Exec(getUserID(GetUserName(req)),getUserID(vars["username"]))
	statement.Close()
	database.Close()
    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func UserUnfollowHandler(res http.ResponseWriter, req *http.Request){
	if (IsEmpty(GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	var database, _ = sql.Open("sqlite3", databasepath)
	statement, _ := database.Prepare("delete from follower where who_id=? and whom_id=?")
	statement.Exec(getUserID(GetUserName(req)),getUserID(vars["username"]))
	statement.Close()
	database.Close()

    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func UserpageRoute(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)	

	if (!checkUsername(vars["username"])){
		http.Redirect(res, req, "/", 302)
	}
	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin": !IsEmpty(GetUserName(req)),
		"username" : vars["username"],
		"postSlice": getUserPosts(vars["username"]),
		"posts": postsAmount(getUserPosts(vars["username"])),
		"visitorUsername" : GetUserName(req),
		"visit" : true,
		"alreadyFollow" : checkIfFollowed(vars["username"],GetUserName(req)),
		"sameUser" : vars["username"] == GetUserName(req),
    }); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func UserpageHandler(res http.ResponseWriter, req *http.Request) {

}

func PersonalTimelineRoute(res http.ResponseWriter, req *http.Request) {


	if (IsEmpty(GetUserName(req))){
		 http.Redirect(res, req, "/", 302)
		}

	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin": !IsEmpty(GetUserName(req)),
		"username" : GetUserName(req),
		"postSlice": getUserPosts(GetUserName(req)),
		"posts": postsAmount(getUserPosts(GetUserName(req))),
    }); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}




func PersonalTimelineHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
 
    text := req.FormValue("text")
 
	_text := false
	
	_text = !IsEmpty(text) 


 
    if _text {

		var database, _ = sql.Open("sqlite3", databasepath)
		statement, _ := database.Prepare("INSERT INTO message (author_id, text, pub_date,flagged) values (?, ?, ?, ?)")
		statement.Exec(getUserID(GetUserName(req)),text,getCurrentTime(),0)
		statement.Close()
		database.Close()

	} else {
		
		fmt.Fprintln(res, "Error")
	}

	PersonalTimelineRoute(res,req)
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

func PublicTimelineRoute(res http.ResponseWriter, req *http.Request) {

	if err := templates["publictimeline"].Execute(res, map[string]interface{}{
		"loggedin": !IsEmpty(GetUserName(req)), 
		"postSlice": getAllPosts(),
		"postSliceLength": len(getAllPosts()),
    }); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}


}

func getAllPosts()[]Post{
	var post Post

	var database, _ = sql.Open("sqlite3", databasepath)
	sqlStatement := `SELECT u.username, m.message_id, m.text, m.pub_date, u.image_url FROM message m join user u ON m.author_id = u.user_id order by m.pub_date desc`
	rows, err := database.Query(sqlStatement)
	if err != nil {
		panic(err)
	   }
	
	defer rows.Close()
	database.Close()
	

	var postSlice []Post
	for rows.Next(){
		rows.Scan(&post.Username, &post.PostMessageid, &post.Text, &post.Date, &post.Image)
		postSlice = append(postSlice, post)
	}

	
	return postSlice;
	

}

func postsAmount(posts []Post) bool{

	if (len(posts) > 0){
		return true
	} else {
		return false
	}
}

func getUserPosts(username string)[]Post{
	var post Post


	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("select m.*, u.image_url  from message m JOIN user u on m.author_id = u.user_id where m.flagged = 0 and m.author_id = u.user_id and (u.user_id = ? or	u.user_id in (select whom_id from follower where who_id = ?)) order by m.pub_date desc")

	if err != nil {
		fmt.Printf("%s", err)
	}


	rows, err := query.Query(getUserID(username),getUserID(username))
	defer rows.Close()
	database.Close()
	
	var postSlice []Post
	for rows.Next(){
		rows.Scan(&post.PostMessageid, &post.AuthorId, &post.Text, &post.Date, &post.Flag, &post.Image)
		post.Username = getUsernameFromID(post.AuthorId)
		postSlice = append(postSlice, post)
	}
	
	return postSlice;

}

func PublicTimelineHandler(res http.ResponseWriter, req *http.Request) {

}

func ContactRoute(res http.ResponseWriter, req *http.Request) {

	if err := templates["contact"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func SigninRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["signin"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// for POST
func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
    name := request.FormValue("username")
    pass := request.FormValue("password")
	redirectTarget := "/"
	_Name, _pwd := false, false
	_Name = !IsEmpty(name)
	_pwd = !IsEmpty(pass)

    if _Name && _pwd {
        // Database check for user data!
        _userIsValid := UserIsValid(name, pass)
		
        if _userIsValid {
            SetCookie(name, response)
            redirectTarget = "/"
        } else {
            redirectTarget = "/register"
        }
    }
    http.Redirect(response, request, redirectTarget, 302)
}

func SignupRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["signup"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
func getGravatarHash(g_email string) string {
	g_email = strings.TrimSpace(g_email)
	g_email = strings.ToLower(g_email)
	h := md5.New()
	io.WriteString(h, g_email)
	finalBytes := h.Sum(nil)
	finalString := hex.EncodeToString(finalBytes)
	return finalString
}

// for POST
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
 
    uName := r.FormValue("username")
    email := r.FormValue("email")
    pwd := r.FormValue("pwd")
 
    _uName, _email, _pwd := false, false, false
    _uName = !IsEmpty(uName)
    _email = !IsEmpty(email)
    _pwd = !IsEmpty(pwd)
 
    if _uName && _email && _pwd {

		if (!checkUsername(uName)){
		

			var database, _ = sql.Open("sqlite3", databasepath)
			gravatar_url := "http://www.gravatar.com/avatar/" + getGravatarHash(email)
			statement, _ := database.Prepare("INSERT INTO user (username, email, pw_hash, image_url) values (?, ?, ?, ?)")
			statement.Exec(uName,email,pwd,gravatar_url)
			database.Close()

			fmt.Fprintln(w, "Username for Register : ", uName)
			fmt.Fprintln(w, "Email for Register : ", email)

		} else {
			fmt.Fprintln(w, "User alrady exits")
		}

    } else {
        fmt.Fprintln(w, "This fields can not be blank!")
    }
}

// for POST
func LogoutHandler(response http.ResponseWriter, request *http.Request) {
    ClearCookie(response)
    http.Redirect(response, request, "/", 302)
}

func SetCookie(userName string, response http.ResponseWriter) {
    value := map[string]string{
        "name": userName,
    }
    if encoded, err := cookieHandler.Encode("cookie", value); err == nil {
        cookie := &http.Cookie{
            Name:  "cookie",
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(response, cookie)
    }
}
 
func ClearCookie(response http.ResponseWriter) {
    cookie := &http.Cookie{
        Name:   "cookie",
        Value:  "",
        Path:   "/",
        MaxAge: -1,
    }
    http.SetCookie(response, cookie)
}
 
func GetUserName(request *http.Request) (userName string) {
    if cookie, err := request.Cookie("cookie"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("cookie", cookie.Value, &cookieValue); err == nil {
            userName = cookieValue["name"]
        }
    }
    return userName
}

func loadTemplates() {
	var baseTemplate = "templates/layout/_base.html"
	templates = make(map[string]*template.Template)

	templates["signin"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signin.html"))
	templates["signup"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signup.html"))
	templates["personaltimeline"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/personal_timeline.html"))
	templates["publictimeline"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/public_timeline.html"))
}

 
func checkUsername(username string) bool {
	var output bool


	var database, _ = sql.Open("sqlite3", databasepath)
	// Prepare your query
	query, err := database.Prepare("SELECT user_id FROM user WHERE username= ?")
	database.Close()

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

func validUser(username string, psw string) bool {
	var output bool


	var database, _ = sql.Open("sqlite3", databasepath)
	// Prepare your query
	query, err := database.Prepare("SELECT user_id FROM user WHERE username = ? AND pw_hash = ?")

	database.Close()
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

	return output

}

func getCurrentTime() string{
	dt := time.Now()
	return (dt.Format("15:04:05 02-01-2006"))
}

func getUserID(username string) int{
	var output int
	// Prepare your query

	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("SELECT user_id FROM user WHERE username = ?")
	database.Close()

	if err != nil {
		fmt.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(username).Scan(&output)

	return output

}

func getUsernameFromID(id int) string{
	var output string

	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("SELECT username FROM user WHERE user_id = ?")
	database.Close()

	if err != nil {
		fmt.Printf("%s", err)
	}
	defer query.Close()

	err = query.QueryRow(id).Scan(&output)

	return output

}

func checkIfFollowed(who string, whom string) bool{
	var output bool

	// Prepare your query
	var database, _ = sql.Open("sqlite3", databasepath)
	query, err := database.Prepare("select * from follower where follower.who_id = ? and follower.whom_id = ?")
	database.Close()

	if err != nil {
		fmt.Printf("%s", err)
	}

	defer query.Close()

	err = query.QueryRow(getUserID(whom), getUserID(who)).Scan(&output)

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

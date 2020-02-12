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
)

var cookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
	
var templates map[string]*template.Template

func init() {
	loadTemplates()
}

var database, _ = sql.Open("sqlite3", "./minitwit.db")

func main() {
	statement, _ := database.Prepare("create table if not exists user (user_id integer primary key autoincrement,username string not null,email string not null,pw_hash string not null);")
	statement.Exec()
	statement2, _ := database.Prepare("create table if not exists follower ( who_id integer, whom_id integer);")
	statement2.Exec()
	statement3, _ := database.Prepare("create table if not exists message (message_id integer primary key autoincrement,author_id integer not null,text string not null,pub_date integer,flagged integer);")
	statement3.Exec()

	fmt.Println("Server starting, point your browser to localhost:3011 to start")

	router := mux.NewRouter()

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	router.HandleFunc("/", IndexRoute).Methods("GET")
	router.HandleFunc("/about", AboutRoute).Methods("GET")
	router.HandleFunc("/contact", ContactRoute).Methods("GET")

	router.HandleFunc("/signin", SigninRoute).Methods("GET")
	router.HandleFunc("/signin", LoginHandler).Methods("POST")


	router.HandleFunc("/signup", SignupRoute).Methods("GET")
	router.HandleFunc("/signup", RegisterHandler).Methods("POST")
	
	router.HandleFunc("/personaltimeline", PersonalTimelineRoute).Methods("GET")
	router.HandleFunc("/personaltimeline", PersonalTimelineHandler).Methods("POST")
	

	router.HandleFunc("/signout", LogoutHandler)

	if err := http.ListenAndServe(":3015", router); err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
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

	if err := templates["index"].Execute(res, map[string]interface{}{
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

func PersonalTimelineRoute(res http.ResponseWriter, req *http.Request) {

	if (IsEmpty(GetUserName(req))){
		 http.Redirect(res, req, "/", 302)
		}

	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin": !IsEmpty(GetUserName(req)),
		"username" : GetUserName(req),
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

		var database, _ = sql.Open("sqlite3", "./minitwit.db")
		statement, _ := database.Prepare("INSERT INTO message (author_id, text, pub_date,flagged) values (?, ?, ?, ?)")
		statement.Exec(getUserID(GetUserName(req)),text,getCurrentTime(),0)
		statement.Close()
		database.Close()

	} else {
		
		fmt.Fprintln(res, "Error")
	}

	PersonalTimelineRoute(res,req)
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
            redirectTarget = "/signup"
        }
    }
    http.Redirect(response, request, redirectTarget, 302)
}

func SignupRoute(res http.ResponseWriter, req *http.Request) {
	if err := templates["signup"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}
// for POST
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
 
    uName := r.FormValue("username")
    email := r.FormValue("email")
    pwd := r.FormValue("password")
    confirmPwd := r.FormValue("confirmPassword")
 
    _uName, _email, _pwd, _confirmPwd := false, false, false, false
    _uName = !IsEmpty(uName)
    _email = !IsEmpty(email)
    _pwd = !IsEmpty(pwd)
    _confirmPwd = !IsEmpty(confirmPwd)
 
    if _uName && _email && _pwd && _confirmPwd {

		if (!checkUsername(uName)){

			statement, _ := database.Prepare("INSERT INTO user (username, email, pw_hash) values (?, ?, ?)")
    		statement.Exec(uName,email,pwd)

			fmt.Fprintln(w, "Username for Register : ", uName)
			fmt.Fprintln(w, "Email for Register : ", email)

		} else {
			fmt.Fprintln(w, "User alrady exits", confirmPwd)
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

	templates["index"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/index.html"))
	templates["about"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/about.html"))
	templates["contact"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/contact.html"))
	templates["signin"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signin.html"))
	templates["signup"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signup.html"))
	templates["personaltimeline"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/personal_timeline.html"))
}

 
func checkUsername(username string) bool {
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

func validUser(username string, psw string) bool {
	var output bool

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

	return output

}

func getCurrentTime() string{
	dt := time.Now()
    return ("" + dt.String())
}

func getUserID(username string) int{
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
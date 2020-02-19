package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	helper "../helpers"
	cookies "../cookies"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	structs "../structs"
)
type Post struct {
	Username string
	PostMessageid int
	AuthorId int 
	Text string
	Date string
	Flag int
	Image string
}


type User struct  {
	user_id int
	username string
	email string
	pw_hash string
	image_url string
  }

  type follower struct  {
	who_id int
	whom_id int
  }
  
  type message struct  {
	message_id int
	author_id int
	text string 
	pub_date string
	flagged int
  }

var databasepath = "/tmp/minitwit.db"

func UserFollowHandler(res http.ResponseWriter, req *http.Request){
	if (helper.IsEmpty(helper.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	var database, _ = gorm.Open("sqlite3", databasepath)

	follow := structs.Follower{Who_id: helper.GetUserID(helper.GetUserName(req)), Whom_id: helper.GetUserID(vars["username"])}
	database.NewRecord(follow)
	database.Create(&follow)
    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func UserUnfollowHandler(res http.ResponseWriter, req *http.Request){
	if (helper.IsEmpty(helper.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	db, err := gorm.Open("sqlite3", databasepath)
		if err != nil {
			panic("failed to connect database")
		}
		defer db.Close()

	follow := structs.Follower{}
	db.Where("who_id = ? AND whom_id = ?", helper.GetUserID(helper.GetUserName(req)),helper.GetUserID(vars["username"])).Delete(follow)



    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
 
    uName := r.FormValue("username")
    email := r.FormValue("email")
    pwd := r.FormValue("pwd")
 
    _uName, _email, _pwd := false, false, false
    _uName = !helper.IsEmpty(uName)
    _email = !helper.IsEmpty(email)
    _pwd = !helper.IsEmpty(pwd)
 
    if _uName && _email && _pwd {
		if (!helper.CheckUsernameExists(uName)){
			var database, _ = gorm.Open("sqlite3", databasepath)
			gravatar_url := "http://www.gravatar.com/avatar/" + helper.GetGravatarHash(email)
		
			user := structs.User{Username: uName, Email: email, Pw_hash: helper.HashPassword(pwd), Image_url: gravatar_url}
			database.NewRecord(user)
			database.Create(&user)
			database.Close()

			cookies.SetCookie(uName, w)
			redirectTarget := "/personaltimeline"
			http.Redirect(w, r, redirectTarget, 302)
		} else {
			fmt.Fprintln(w, "User already exits")
		}

    } else {
        fmt.Fprintln(w, "This fields can not be blank!")
    }
}

func LoginHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
    name := request.FormValue("username")
    pass := request.FormValue("password")
	redirectTarget := "/"
	_Name, _pwd := false, false
	_Name = !helper.IsEmpty(name)
	_pwd = !helper.IsEmpty(pass)

    if _Name && _pwd {
        // Database check for user data!
        _userIsValid := helper.UserIsValid(name, pass)
		
        if _userIsValid {
            cookies.SetCookie(name, response)
            redirectTarget = "/"
        } else {
            redirectTarget = "/register"
        }
    }
    http.Redirect(response, request, redirectTarget, 302)
}

func PersonalTimelineHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()

    text := req.FormValue("text")
	_text := false
	_text = !helper.IsEmpty(text) 

    if _text {
		var database, _ = gorm.Open("sqlite3", databasepath)
		message := structs.Message{Author_id:helper.GetUserID(helper.GetUserName(req)),Text:text,Pub_date:helper.GetCurrentTime(),Flagged:0}
		database.NewRecord(message)
		database.Create(&message)
		database.Close()

	} else {	
		fmt.Fprintln(res, "Error")
	}

	PersonalTimelineRoute(res,req)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
    cookies.ClearCookie(response)
    http.Redirect(response, request, "/", 302)
}


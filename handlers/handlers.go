package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"database/sql"
	
	"fmt"

	helpers "../helpers"
	cookies "../cookies"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	structs "../structs"
)

var db *sql.DB

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



func UserFollowHandler(res http.ResponseWriter, req *http.Request){
	if (helpers.IsEmpty(helpers.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	db := helpers.GetDB()

	follow := structs.Follower{Who_id: helpers.GetUserID(helpers.GetUserName(req)), Whom_id: helpers.GetUserID(vars["username"])}
	db.NewRecord(follow)
	db.Create(&follow)
    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}




func UserUnfollowHandler(res http.ResponseWriter, req *http.Request){
	if (helpers.IsEmpty(helpers.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	db := helpers.GetDB()
	defer db.Close()

	follow := structs.Follower{}
	db.Where("who_id = ? AND whom_id = ?", helpers.GetUserID(helpers.GetUserName(req)),helpers.GetUserID(vars["username"])).Delete(follow)



    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
 
    uName := r.FormValue("username")
    email := r.FormValue("email")
    pwd := r.FormValue("pwd")
 
    _uName, _email, _pwd := false, false, false
    _uName = !helpers.IsEmpty(uName)
    _email = !helpers.IsEmpty(email)
    _pwd = !helpers.IsEmpty(pwd)
 
    if _uName && _email && _pwd {
		if (!helpers.CheckUsernameExists(uName)){
			db := helpers.GetDB()
			gravatar_url := "http://www.gravatar.com/avatar/" + helpers.GetGravatarHash(email)
		
			user := structs.User{Username: uName, Email: email, Pw_hash: helpers.HashPassword(pwd), Image_url: gravatar_url}
			db.NewRecord(user)
			db.Create(&user)
			db.Close()

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
	_Name = !helpers.IsEmpty(name)
	_pwd = !helpers.IsEmpty(pass)

    if _Name && _pwd {
        // Database check for user data!
        _userIsValid := helpers.UserIsValid(name, pass)
		
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
	_text = !helpers.IsEmpty(text) 

    if _text {
		db := helpers.GetDB()
		message := structs.Message{Author_id:helpers.GetUserID(helpers.GetUserName(req)),Text:text,Pub_date:helpers.GetCurrentTime(),Flagged:0}
		db.NewRecord(message)
		db.Create(&message)
		db.Close()

	} else {	
		fmt.Fprintln(res, "Error")
	}

	PersonalTimelineRoute(res,req)
}

func LogoutHandler(response http.ResponseWriter, request *http.Request) {
    cookies.ClearCookie(response)
    http.Redirect(response, request, "/", 302)
}


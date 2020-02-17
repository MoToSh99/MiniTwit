package handlers

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"database/sql"
	helper "../helpers"
	cookies "../cookies"
)


var databasepath = "/tmp/minitwit.db"

func UserFollowHandler(res http.ResponseWriter, req *http.Request){
	if (helper.IsEmpty(helper.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	var database, _ = sql.Open("sqlite3", databasepath)
	statement, _ := database.Prepare("insert into follower (who_id, whom_id) values (?, ?)")
	statement.Exec(helper.GetUserID(helper.GetUserName(req)),helper.GetUserID(vars["username"]))
	statement.Close()
	database.Close()
    http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func UserUnfollowHandler(res http.ResponseWriter, req *http.Request){
	if (helper.IsEmpty(helper.GetUserName(req))){
		res.WriteHeader(http.StatusUnauthorized)
	}
	
	vars := mux.Vars(req)

	var database, _ = sql.Open("sqlite3", databasepath)
	statement, _ := database.Prepare("delete from follower where who_id=? and whom_id=?")
	statement.Exec(helper.GetUserID(helper.GetUserName(req)),helper.GetUserID(vars["username"]))
	statement.Close()
	database.Close()

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
		

			var database, _ = sql.Open("sqlite3", databasepath)
			gravatar_url := "http://www.gravatar.com/avatar/" + helper.GetGravatarHash(email)
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

		var database, _ = sql.Open("sqlite3", databasepath)
		statement, _ := database.Prepare("INSERT INTO message (author_id, text, pub_date,flagged) values (?, ?, ?, ?)")
		statement.Exec(helper.GetUserID(helper.GetUserName(req)),text,helper.GetCurrentTime(),0)
		statement.Close()
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


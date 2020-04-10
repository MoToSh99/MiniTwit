package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	cookies "../cookies"
	helpers "../helpers"
	structs "../structs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var db *gorm.DB

type Post struct {
	Username   string
	Message_id int
	Author_id  int
	Text       string
	Pub_date   string
	Flagged    int
	Image_url  string
}

type User struct {
	user_id   int
	username  string
	email     string
	pw_hash   string
	image_url string
}

type follower struct {
	who_id  int
	whom_id int
}

type message struct {
	message_id int
	author_id  int
	text       string
	pub_date   string
	flagged    int
}

func AddSafeHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Strict-Transport-Security", "max-age=2592000; includeSubDomains")
}

func UserFollowHandler(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if helpers.IsEmpty(helpers.GetUserName(req)) {

		res.WriteHeader(http.StatusUnauthorized)
	}

	vars := mux.Vars(req)

	db := helpers.GetDB()

	follow := structs.Follower{Who_id: helpers.GetUserID(helpers.GetUserName(req)), Whom_id: helpers.GetUserID(vars["username"])}
	db.NewRecord(follow)
	db.Create(&follow)
	http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)
}

func UserUnfollowHandler(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if helpers.IsEmpty(helpers.GetUserName(req)) {
		res.WriteHeader(http.StatusUnauthorized)
	}

	vars := mux.Vars(req)

	db := helpers.GetDB()

	follow := structs.Follower{}
	db.Where("who_id = ? AND whom_id = ?", helpers.GetUserID(helpers.GetUserName(req)), helpers.GetUserID(vars["username"])).Delete(follow)

	http.Redirect(res, req, fmt.Sprintf("/%v", vars["username"]), 302)

}

func RegisterHandler(res http.ResponseWriter, r *http.Request) {
	AddSafeHeaders(res)
	r.ParseForm()

	uName := r.FormValue("username")
	email := r.FormValue("email")
	pwd := r.FormValue("pwd")
	pwdconfirm := r.FormValue("confirmPassword")

	_uName, _email, _pwd, _pwdconfirm := false, false, false, false
	_uName = !helpers.IsEmpty(uName)
	_email = !helpers.IsEmpty(email)
	_pwd = !helpers.IsEmpty(pwd)
	_pwdconfirm = !helpers.IsEmpty(pwdconfirm)

	errorMsg := ""

	if !_uName {
		errorMsg = "You have to enter a username"
	} else if !_email {
		errorMsg = "You have to enter a valid email address"
	} else if !_pwd {
		errorMsg = "You have to enter a password"
	} else if !_pwdconfirm {
		errorMsg = "You have to confirm your password"
	} else if pwd != pwdconfirm {
		errorMsg = "The two passwords do not match"
	} else if helpers.CheckUsernameExists(uName) {
		errorMsg = "The username is already taken"
	} else {
		db := helpers.GetDB()
		gravatar_url := "http://www.gravatar.com/avatar/" + helpers.GetGravatarHash(email)

		user := structs.User{Username: uName, Email: email, Pw_hash: helpers.HashPassword(pwd), Image_url: gravatar_url}
		db.NewRecord(user)
		db.Create(&user)

		cookies.SetCookie(uName, res)
		redirectTarget := "/personaltimeline"
		http.Redirect(res, r, redirectTarget, 302)
	}
	ShowSignUpError(res, r, errorMsg)

}

func ShowSignUpError(res http.ResponseWriter, r *http.Request, errorMsg string) {
	AddSafeHeaders(res)
	if err := templates["signup"].Execute(res, map[string]interface{}{
		"error":           true,
		"FlashedMessages": errorMsg,
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func LoginHandler(res http.ResponseWriter, request *http.Request) {
	AddSafeHeaders(res)
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
			cookies.SetCookie(name, res)
			redirectTarget = "/"
		} else {
			redirectTarget = "/register"
		}
	}
	http.Redirect(res, request, redirectTarget, 302)
}

func PersonalTimelineHandler(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	req.ParseForm()

	text := req.FormValue("text")
	_text := false
	_text = !helpers.IsEmpty(text)

	if _text {
		db := helpers.GetDB()
		message := structs.Message{Author_id: helpers.GetUserID(helpers.GetUserName(req)), Text: text, Pub_date: helpers.GetCurrentTime(), Flagged: 0}
		db.NewRecord(message)
		db.Create(&message)

	} else {
		fmt.Fprintln(res, "Error")
	}

	PersonalTimelineRoute(res, req)
}

func LogoutHandler(res http.ResponseWriter, request *http.Request) {
	AddSafeHeaders(res)
	cookies.ClearCookie(res)
	http.Redirect(res, request, "/", 302)
}

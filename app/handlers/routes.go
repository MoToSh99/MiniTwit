package handlers

import (
	"html/template"
	"net/http"

	helper "../helpers"
	"github.com/gorilla/mux"
)

var templates map[string]*template.Template

func AddSafeHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("Strict-Transport-Security", "max-age=2592000; includeSubDomains")
}

func ContactRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["contact"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func SigninRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["signin"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func SignupRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["signup"].Execute(res, nil); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func PublicTimelineRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["publictimeline"].Execute(res, map[string]interface{}{
		"loggedin":        !helper.IsEmpty(helper.GetUserName(req)),
		"postSlice":       helper.GetMoreposts(10),
		"postSliceLength": len(helper.GetMoreposts(10)),
		"showMoreActive":  true,
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func PublicTimelineLoadMore(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["publictimeline"].Execute(res, map[string]interface{}{
		"loggedin":        !helper.IsEmpty(helper.GetUserName(req)),
		"postSlice":       helper.GetAllPosts(),
		"postSliceLength": len(helper.GetAllPosts()),
		"showMoreActive":  false,
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func PersonalTimelineRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if helper.IsEmpty(helper.GetUserName(req)) {
		http.Redirect(res, req, "/", 302)
	}

	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin":  !helper.IsEmpty(helper.GetUserName(req)),
		"username":  helper.GetUserName(req),
		"postSlice": helper.GetUserPosts(helper.GetUserName(req)),
		"posts":     helper.PostsAmount(helper.GetUserPosts(helper.GetUserName(req))),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func PersonalTimelineMore(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if helper.IsEmpty(helper.GetUserName(req)) {
		http.Redirect(res, req, "/", 302)
	}

	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin":  !helper.IsEmpty(helper.GetUserName(req)),
		"username":  helper.GetUserName(req),
		"postSlice": helper.GetUserPosts(helper.GetUserName(req)),
		"posts":     helper.PostsAmount(helper.GetUserPosts(helper.GetUserName(req))),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func UserpageRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	vars := mux.Vars(req)

	if !helper.CheckUsernameExists(vars["username"]) {
		http.Redirect(res, req, "/", 302)
	}
	if err := templates["personaltimeline"].Execute(res, map[string]interface{}{
		"loggedin":        !helper.IsEmpty(helper.GetUserName(req)),
		"username":        vars["username"],
		"postSlice":       helper.GetUserPosts(vars["username"]),
		"posts":           helper.PostsAmount(helper.GetUserPosts(vars["username"])),
		"visitorUsername": helper.GetUserName(req),
		"visit":           true,
		"alreadyFollow":   helper.CheckIfFollowed(vars["username"], helper.GetUserName(req)),
		"sameUser":        vars["username"] == helper.GetUserName(req),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func IndexRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["publictimeline"].Execute(res, map[string]interface{}{
		"loggedin": !helper.IsEmpty(helper.GetUserName(req)),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func AboutRoute(res http.ResponseWriter, req *http.Request) {
	AddSafeHeaders(res)
	if err := templates["about"].Execute(res, map[string]interface{}{
		"loggedin": !helper.IsEmpty(helper.GetUserName(req)),
	}); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func LoadTemplates() {
	var baseTemplate = "templates/layout/_base.html"
	templates = make(map[string]*template.Template)

	templates["signin"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signin.html"))
	templates["signup"] = template.Must(template.ParseFiles(baseTemplate, "templates/account/signup.html"))
	templates["personaltimeline"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/personal_timeline.html"))
	templates["publictimeline"] = template.Must(template.ParseFiles(baseTemplate, "templates/home/public_timeline.html"))
}

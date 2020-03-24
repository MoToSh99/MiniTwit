package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	helpers "../helpers"
	metrics "../metrics"
	structs "../structs"
	"github.com/Jeffail/gabs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

var gLATEST = 0

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Pwd      string `json:"pwd"`
}

type Message struct {
	Username string `json:"username"`
	Text     string `json:"content"`
}

func Not_req_from_simulator(w http.ResponseWriter, r *http.Request) {
	from_simulator := r.Header.Get("Authorization")
	if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		error := "You are not authorized to use this resource!"
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(fmt.Sprintf(`{"status": 403, "error_msg": %v}`, error)))
	}
}

func Update_latest(res http.ResponseWriter, req *http.Request) {
	jsonint, _ := strconv.Atoi(req.URL.Query().Get("latest"))
	if jsonint != 0 {
		gLATEST = jsonint
	}
}

// Get_latest godoc
// @Summary Get latest accepted id
// @Produce json
// @Success 200
// @Router /latest [get]
func Get_latest(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"latest": %v}`, gLATEST)))

	elapsed := time.Since(start)
	metrics.ResponseTimeRegister.Observe(float64(elapsed.Seconds()*1000))
}


// Register godoc
// @Summary Post new user to register
// @Produce json
// @Param name string true "User Name"
// @Success 200
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	Update_latest(w, r)

	var user User

	var error = ""
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if helpers.IsEmpty(user.Username) {
		error = "You have to enter a username"
	} else if helpers.IsEmpty(user.Email) || !strings.Contains(user.Email, "@") {
		error = "You have to enter a valid email address"
	} else if helpers.IsEmpty(user.Pwd) {
		error = "You have to enter a password"
	} else if helpers.CheckUsernameExists(user.Username) {
		error = "The username is already taken"
		return
	} else {

		db := helpers.GetDB()

		gravatar_url := "http://www.gravatar.com/avatar/" + helpers.GetGravatarHash(user.Email) + "?&d=identicon"

		metrics.UsersRegistered.Inc()

		db.Create(&structs.User{Username: user.Username, Email: user.Email, Pw_hash: helpers.HashPassword(user.Pwd), Image_url: gravatar_url})

	}

	if !helpers.IsEmpty(error) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"status": 400, "error_msg": %v}`, error)))
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
	elapsed := time.Since(start)
	metrics.ResponseTimeRegister.Observe(float64(elapsed.Seconds()*1000))
}

type Post struct {
	Text     string `json:"content"`
	Pub_date string `json:"pub_date"`
	Username string `json:"user"`
}

func Messages(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	Update_latest(w, r)
	Not_req_from_simulator(w, r)

	no, _ := strconv.Atoi(r.URL.Query().Get("no"))
	db := helpers.GetDB()

	var postSlice []Post

	db.Table("messages").Limit(no).Order("messages.pub_date").Select("messages.text, messages.pub_date, users.username").Joins("join users on users.user_id = messages.author_id").Scan(&postSlice)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(postSlice)

	elapsed := time.Since(start)
	metrics.ResponseTimeMsgs.Observe(float64(elapsed.Seconds()*1000))
}

func Messages_per_user(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	Update_latest(w, r)
	Not_req_from_simulator(w, r)
	vars := mux.Vars(r)

	if !helpers.CheckUsernameExists(vars["username"]) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		db := helpers.GetDB()
		no, _ := strconv.Atoi(r.URL.Query().Get("no"))

		var postSlice []Post

		db.Table("messages").Limit(no).Order("messages.pub_date").Select("messages.text, messages.pub_date, users.username").Joins("join users on users.user_id = messages.author_id").Where("messages.flagged = 0 AND users.username = ?", vars["username"]).Scan(&postSlice)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(postSlice)

	} else if r.Method == http.MethodPost {

		var msg Message

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db := helpers.GetDB()

		db.Create(&structs.Message{Author_id: helpers.GetUserID(vars["username"]), Text: msg.Text, Pub_date: helpers.GetCurrentTime(), Flagged: 0})

		metrics.MessagesSent.Inc()

		w.WriteHeader(http.StatusNoContent)
	}

	elapsed := time.Since(start)
	metrics.ResponseTimeMsgsPerUser.Observe(float64(elapsed.Seconds()*1000))
}

type FollowUser struct {
	Follows_username  string `json:"follow"`
	Unfollow_username string `json:"unfollow"`
}

func Follow(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	Update_latest(w, r)
	Not_req_from_simulator(w, r)
	vars := mux.Vars(r)
	no, _ := strconv.Atoi(r.URL.Query().Get("no"))
	var follow FollowUser

	if r.Body != http.NoBody {
		err := json.NewDecoder(r.Body).Decode(&follow)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if !helpers.CheckUsernameExists(vars["username"]) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost && !helpers.IsEmpty(follow.Follows_username) {
		follows_username := follow.Follows_username

		if !helpers.CheckUsernameExists(follows_username) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		db := helpers.GetDB()
		metrics.UsersFollowed.Inc()
		db.Create(&structs.Follower{Who_id: helpers.GetUserID(vars["username"]), Whom_id: helpers.GetUserID(follows_username)})

	} else if r.Method == http.MethodPost && !helpers.IsEmpty(follow.Unfollow_username) {
		unfollows_username := follow.Unfollow_username
		if !helpers.CheckUsernameExists(unfollows_username) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if !helpers.CheckUsernameExists(unfollows_username) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		db := helpers.GetDB()

		follow := structs.Follower{}

		metrics.UsersUnfollowed.Inc()

		db.Where("who_id = ? AND whom_id = ?", helpers.GetUserID(vars["username"]), helpers.GetUserID(unfollows_username)).Delete(follow)

	} else if r.Method == http.MethodGet {
		db := helpers.GetDB()

		userSlice := []structs.Follower{}

		db.Limit(no).Where("who_id = ?", helpers.GetUserID(vars["username"])).Order("whom_id").Find(&userSlice)

		jsonObj := gabs.New()
		jsonObj.Array("follows")
		for _, v := range userSlice {
			jsonObj.ArrayAppend(helpers.GetUsernameFromID(v.Whom_id), "follows")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, jsonObj.StringIndent("", "  "))

	}
	elapsed := time.Since(start)
	metrics.ResponseTimeFollow.Observe(float64(elapsed.Seconds()*1000))
}

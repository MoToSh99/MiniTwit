package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"
	helpers "../helpers"
	logger "../logger"
	metrics "../metrics"
	structs "../structs"
	handlers "../handlers"
	"github.com/Jeffail/gabs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func GetConnString() string {
	var connString = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USERNAME"), os.Getenv("DB_DATABASE"), os.Getenv("DB_PASSWORD"))

	return connString
}


func GetDB() *gorm.DB {

	var connString = GetConnString()

	db, err := gorm.Open("postgres", connString)


	if err != nil {
		panic("failed to connect database")
	}

	return db

}

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

func Not_req_from_simulator(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	from_simulator := r.Header.Get("Authorization")
	if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh" {
		error := "You are not authorized to use this resource!"
		res.WriteHeader(http.StatusForbidden)
		res.Write([]byte(fmt.Sprintf(`{"status": 403, "error_msg": %v}`, error)))
	}
}

func Update_latest(res http.ResponseWriter, req *http.Request) {
	handlers.AddSafeHeaders(res)
	jsonint, _ := strconv.Atoi(req.URL.Query().Get("latest"))
	if jsonint != 0 {
		gLATEST = jsonint
	}
}

// Get_latest godoc
// @Summary Get latest accepted id
// @Produce json
// @Success 200 "Returns latest accepted id by api"
// @Router /latest [get]
func Get_latest(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	start := time.Now()

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(fmt.Sprintf(`{"latest": %v}`, gLATEST)))
	logger.Send(fmt.Sprintf(`API - Sent latest value:  %v`, gLATEST))
	elapsed := time.Since(start)
	metrics.ResponseTimeRegister.Observe(float64(elapsed.Seconds() * 1000))
}

// Register godoc
// @Summary Post new user to register
// @Produce json
// @Param name path string true "User Name"
// @Param email path string true "Email"
// @Param password path string true "Password"
// @Success 204 "User registered"
// @Failure 400 "Error on insert with description"
// @Router /register [post]
func Register(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	start := time.Now()

	Update_latest(res, r)

	var user User

	var error = ""
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
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

		db := GetDB()

		gravatar_url := "http://www.gravatar.com/avatar/" + helpers.GetGravatarHash(user.Email) + "?&d=identicon"

		metrics.UsersRegistered.Inc()

		db.Create(&structs.User{Username: user.Username, Email: user.Email, Pw_hash: helpers.HashPassword(user.Pwd), Image_url: gravatar_url})

		logger.Send(fmt.Sprintf(`API - User registered with username:  %v`, user.Username))

	}

	if !helpers.IsEmpty(error) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf(`{"status": 400, "error_msg": %v}`, error)))
	} else {
		res.WriteHeader(http.StatusNoContent)
	}
	elapsed := time.Since(start)
	metrics.ResponseTimeRegister.Observe(float64(elapsed.Seconds() * 1000))
}

type Post struct {
	Text     string `json:"content"`
	Pub_date string `json:"pub_date"`
	Username string `json:"user"`
}

func Messages(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	start := time.Now()
	Update_latest(res, r)
	Not_req_from_simulator(res, r)

	no, _ := strconv.Atoi(r.URL.Query().Get("no"))
	db := GetDB()

	var postSlice []Post

	db.Table("messages").Limit(no).Order("messages.pub_date").Select("messages.text, messages.pub_date, users.username").Joins("join users on users.user_id = messages.author_id").Scan(&postSlice)

	logger.Send(fmt.Sprintf(`API - Sent list of messages wit limit:  %v`, no))

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(postSlice)

	elapsed := time.Since(start)
	metrics.ResponseTimeMsgs.Observe(float64(elapsed.Seconds() * 1000))
}

func Messages_per_user(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	start := time.Now()
	Update_latest(res, r)
	Not_req_from_simulator(res, r)
	vars := mux.Vars(r)

	if !helpers.CheckUsernameExists(vars["username"]) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodGet {
		db := GetDB()
		no, _ := strconv.Atoi(r.URL.Query().Get("no"))

		var postSlice []Post

		db.Table("messages").Limit(no).Order("messages.pub_date").Select("messages.text, messages.pub_date, users.username").Joins("join users on users.user_id = messages.author_id").Where("messages.flagged = 0 AND users.username = ?", vars["username"]).Scan(&postSlice)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(postSlice)

		logger.Send(fmt.Sprintf(`API - Sent list of messages for user:  %v`, vars["username"]))

	} else if r.Method == http.MethodPost {

		var msg Message

		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		db := GetDB()

		db.Create(&structs.Message{Author_id: helpers.GetUserID(vars["username"]), Text: msg.Text, Pub_date: helpers.GetCurrentTime(), Flagged: 0})

		logger.Send(fmt.Sprintf(`API - Posted message with username:  %v`, vars["username"]))

		metrics.MessagesSent.Inc()

		res.WriteHeader(http.StatusNoContent)
	}

	elapsed := time.Since(start)
	metrics.ResponseTimeMsgsPerUser.Observe(float64(elapsed.Seconds() * 1000))
}

type FollowUser struct {
	Follows_username  string `json:"follow"`
	Unfollow_username string `json:"unfollow"`
}

func Follow(res http.ResponseWriter, r *http.Request) {
	handlers.AddSafeHeaders(res)
	start := time.Now()
	Update_latest(res, r)
	Not_req_from_simulator(res, r)
	vars := mux.Vars(r)
	no, _ := strconv.Atoi(r.URL.Query().Get("no"))
	var follow FollowUser

	if r.Body != http.NoBody {
		err := json.NewDecoder(r.Body).Decode(&follow)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if !helpers.CheckUsernameExists(vars["username"]) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost && !helpers.IsEmpty(follow.Follows_username) {
		follows_username := follow.Follows_username

		if !helpers.CheckUsernameExists(follows_username) {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		db := GetDB()
		metrics.UsersFollowed.Inc()
		db.Create(&structs.Follower{Who_id: helpers.GetUserID(vars["username"]), Whom_id: helpers.GetUserID(follows_username)})

		logger.Send(fmt.Sprintf(`API - %v follows %v`, vars["username"], follows_username))

	} else if r.Method == http.MethodPost && !helpers.IsEmpty(follow.Unfollow_username) {
		unfollows_username := follow.Unfollow_username
		if !helpers.CheckUsernameExists(unfollows_username) {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		if !helpers.CheckUsernameExists(unfollows_username) {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		db := GetDB()

		follow := structs.Follower{}

		metrics.UsersUnfollowed.Inc()

		db.Where("who_id = ? AND whom_id = ?", helpers.GetUserID(vars["username"]), helpers.GetUserID(unfollows_username)).Delete(follow)

		logger.Send(fmt.Sprintf(`API - %v unfollows %v `, vars["username"], unfollows_username))

	} else if r.Method == http.MethodGet {
		db := GetDB()

		userSlice := []structs.Follower{}

		db.Limit(no).Where("who_id = ?", helpers.GetUserID(vars["username"])).Order("whom_id").Find(&userSlice)

		jsonObj := gabs.New()
		jsonObj.Array("follows")
		for _, v := range userSlice {
			jsonObj.ArrayAppend(helpers.GetUsernameFromID(v.Whom_id), "follows")
		}

		logger.Send(fmt.Sprintf(`API - List of followers sent for username: %v `, vars["username"]))

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		fmt.Fprintln(res, jsonObj.StringIndent("", "  "))

	}
	elapsed := time.Since(start)
	metrics.ResponseTimeFollow.Observe(float64(elapsed.Seconds() * 1000))
}

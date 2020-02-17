package api

import (
    "fmt"
    "net/http"
    "encoding/json"
    helpers "../helpers"
    "strconv"
    "strings"
    "database/sql"
    "github.com/gorilla/mux"
    "github.com/Jeffail/gabs"
)
var databasepath = "/tmp/minitwit.db"

var gLATEST = 0;

type User struct {
	Username string `json:"username"`
    Email  string `json:"email"`
    Pwd string `json:"pwd"`
}

type Message struct {
	Username string `json:"username"`
    Text  string `json:"content"`
}




func Test(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(res, "You have arrived at the API Server!")
}


func Not_req_from_simulator(w http.ResponseWriter, r *http.Request){
    from_simulator := r.Header.Get("Authorization")
    if from_simulator != "Basic c2ltdWxhdG9yOnN1cGVyX3NhZmUh"{
        error := "You are not authorized to use this resource!"
        w.WriteHeader(http.StatusForbidden)
        w.Write([]byte(fmt.Sprintf(`{"status": 403, "error_msg": %v}`, error)))
    }

    
}


func Update_latest(res http.ResponseWriter, req *http.Request){

    jsonint, _ := strconv.Atoi(req.URL.Query().Get("latest"))

    if (jsonint != 0){

        gLATEST = jsonint
    }
}



func Get_latest(w http.ResponseWriter, r *http.Request){
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write([]byte(fmt.Sprintf(`{"latest": %v}`, gLATEST)))

}


func Register(w http.ResponseWriter, r *http.Request){
    Update_latest(w,r)

    var user User

    var error = ""
   
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if (helpers.IsEmpty(user.Username)){
        error = "You have to enter a username"
    } else if (helpers.IsEmpty(user.Email) || !strings.Contains(user.Email, "@") ){
        error = "You have to enter a valid email address"
    } else if (helpers.IsEmpty(user.Pwd) ){
        error = "You have to enter a password"
    } else if (helpers.CheckUsernameExists(user.Username) ){
        error = "The username is already taken"
    } else{
        var database, _ = sql.Open("sqlite3", databasepath)
        gravatar_url := "http://www.gravatar.com/avatar/" + helpers.GetGravatarHash(user.Email)
		statement, _ := database.Prepare("INSERT INTO user (username, email, pw_hash, image_url) values (?, ?, ?, ?)")
        statement.Exec(user.Username,user.Email,user.Pwd,gravatar_url )
        database.Close()
    }

    if (!helpers.IsEmpty(error)){
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf(`{"status": 400, "error_msg": %v}`, error)))
    } else{
        w.WriteHeader(http.StatusNoContent)
    }
}


type Post struct {
    Text string `json:"content"`
	Date string `json:"pub_Date"`
	Username string `json:"user"`
}


func Messages(w http.ResponseWriter, r *http.Request){
    Update_latest(w,r)
   Not_req_from_simulator(w,r)


    var database, _ = sql.Open("sqlite3", databasepath)
    query, _ := database.Prepare("SELECT m.text, m.pub_date, u.username FROM message m, user u WHERE m.flagged = 0 AND m.author_id = u.user_id ORDER BY m.pub_date DESC LIMIT ?")
    no, _ := strconv.Atoi(r.URL.Query().Get("no"))
    rows, _ := query.Query(no)
    database.Close()


    var post Post
    var postSlice []Post

    for rows.Next() {
        rows.Scan(&post.Text, &post.Date, &post.Username)
        postSlice = append(postSlice, post)
        

    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(postSlice)
    

}



func Messages_per_user(w http.ResponseWriter, r *http.Request){
    Update_latest(w,r)
    Not_req_from_simulator(w,r)
    vars := mux.Vars(r)	
    var post Post

    if (!helpers.CheckUsernameExists(vars["username"])){
        w.WriteHeader(http.StatusNotFound)
    }

    if (r.Method == http.MethodGet) {
        var database, _ = sql.Open("sqlite3", databasepath)
        query, _ := database.Prepare("SELECT m.text, m.pub_date, u.username FROM message m, user u WHERE m.flagged = 0 AND m.author_id = u.user_id AND u.user_id = ?ORDER BY m.pub_date DESC LIMIT ?")
        no, _ := strconv.Atoi(r.URL.Query().Get("no"))
        rows, _ := query.Query(helpers.GetUserID(vars["username"]),no)
        database.Close()
        
        var postSlice []Post

        for rows.Next() {
            rows.Scan(&post.Text, &post.Date, &post.Username)
            postSlice = append(postSlice, post)
            

        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(postSlice)

    } else if (r.Method == http.MethodPost) {
        
        var msg Message

        err := json.NewDecoder(r.Body).Decode(&msg)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        
		var database, _ = sql.Open("sqlite3", databasepath)
		statement, _ := database.Prepare("INSERT INTO message (author_id, text, pub_date,flagged) values (?, ?, ?, ?)")
        statement.Exec(helpers.GetUserID(vars["username"]),msg.Text ,helpers.GetCurrentTime(),0)
        statement.Close()
        database.Close()
        w.WriteHeader(http.StatusNoContent)
        
    }

}

type FollowUser struct {
	Follows_username string `json:"follow"`
	Unfollow_username string `json:"unfollow"`
}



func Follow(w http.ResponseWriter, r *http.Request){
    Update_latest(w,r)
    Not_req_from_simulator(w,r)
    vars := mux.Vars(r)	

    var follow FollowUser

    if r.Body != http.NoBody {
        err := json.NewDecoder(r.Body).Decode(&follow)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    }


    if (!helpers.CheckUsernameExists(vars["username"])){
        w.WriteHeader(http.StatusNotFound)
    }

    if (r.Method == http.MethodPost &&  !helpers.IsEmpty(follow.Follows_username)) {
        follows_username := follow.Follows_username
        if (!helpers.CheckUsernameExists(follows_username)){
            w.WriteHeader(http.StatusNotFound)
        }
        var database, _ = sql.Open("sqlite3", databasepath)
        statement, _ := database.Prepare("insert into follower (who_id, whom_id) values (?, ?)")
        statement.Exec(helpers.GetUserID(vars["username"]),helpers.GetUserID(follows_username))
        statement.Close()
        database.Close()

    } else if (r.Method == http.MethodPost &&  !helpers.IsEmpty(follow.Unfollow_username)){
        unfollows_username := follow.Unfollow_username
        if (!helpers.CheckUsernameExists(unfollows_username)){
            w.WriteHeader(http.StatusNotFound)
        }
        var database, _ = sql.Open("sqlite3", databasepath)
        statement, _ := database.Prepare("delete from follower where who_id=? and whom_id=?")
        statement.Exec(helpers.GetUserID(vars["username"]),helpers.GetUserID(unfollows_username))
        statement.Close()
        database.Close()

    } else if (r.Method == http.MethodGet){
        var database, _ = sql.Open("sqlite3", databasepath)
        query, _ := database.Prepare("SELECT user.username FROM user INNER JOIN follower ON follower.whom_id=user.user_id WHERE follower.who_id=? LIMIT ?")
        no, _ := strconv.Atoi(r.URL.Query().Get("no"))
        rows, _ := query.Query(helpers.GetUserID(vars["username"]),no)
        database.Close()
    
       
        var user string

        jsonObj := gabs.New()
        jsonObj.Array("follows")
        for rows.Next() {
            rows.Scan(&user)
            jsonObj.ArrayAppend(user, "follows")
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w,jsonObj.StringIndent("", "  "))

        
    }
}

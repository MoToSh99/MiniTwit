
type Post struct {
    Text string `json:"content"`
	Date string `json:"pub_Date"`
	Username string `json:"user"`
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
        query, _ := database.Prepare("SELECT m.text, m.pub_date, u.username FROM message m, user u WHERE m.flagged = 0 AND m.author_id = u.user_id AND u.user_id = ?ORDER BY m.pub_date DESC LIMIT ?")
        no, _ := strconv.Atoi(r.URL.Query().Get("no"))
        rows, _ := query.Query(helpers.GetUserID(vars["username"]),no)
        
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
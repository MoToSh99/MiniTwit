package main

import (
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
	api "./api"
	handler "./handlers"
	structs "./structs"
)


func init() {
	handler.LoadTemplates()
}

var databasepath = "/tmp/minitwit.db"

func main() {
	db, err := gorm.Open("sqlite3", databasepath)
	
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&structs.User{}, &structs.Follower{}, &structs.Message{})

	router := mux.NewRouter()


	router.HandleFunc("/favicon.ico", faviconHandler)

	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	router.HandleFunc("/", handler.PublicTimelineRoute).Methods("GET")
	router.HandleFunc("/about", handler.AboutRoute).Methods("GET")
	router.HandleFunc("/contact", handler.ContactRoute).Methods("GET")

	router.HandleFunc("/signin", handler.SigninRoute).Methods("GET")
	router.HandleFunc("/signin", handler.LoginHandler).Methods("POST")


	router.HandleFunc("/register", handler.SignupRoute).Methods("GET")
	router.HandleFunc("/register", handler.RegisterHandler).Methods("POST")
	
	router.HandleFunc("/personaltimeline", handler.PersonalTimelineRoute).Methods("GET")
	router.HandleFunc("/personaltimeline", handler.PersonalTimelineHandler).Methods("POST")
	

	router.HandleFunc("/signout", handler.LogoutHandler)

	router.HandleFunc("/publictimeline", handler.PublicTimelineRoute).Methods("GET")
	

	router.HandleFunc("/{username}", handler.UserpageRoute).Methods("GET")

	router.HandleFunc("/{username}/follow", handler.UserFollowHandler)
	router.HandleFunc("/{username}/unfollow", handler.UserUnfollowHandler)

	

	apiRoute := mux.NewRouter()
	//apiRoute.HandleFunc("/test", api.Test)
	apiRoute.HandleFunc("/latest", api.Get_latest)
	apiRoute.HandleFunc("/register", api.Register).Methods("POST")
	apiRoute.HandleFunc("/msgs", api.Messages)
	apiRoute.HandleFunc("/msgs/{username}", api.Messages_per_user)
	apiRoute.HandleFunc("/fllws/{username}", api.Follow)




	port := 5000
	log.Printf("Server starting on port %v\n", port)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))}()

	apiport := 5001
	log.Printf("Api Server starting on port %v\n", apiport)
    go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", apiport), apiRoute))}()
	
    select {}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/public/favicon.ico")
}


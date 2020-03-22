package main

import (
	"fmt"
	"log"
	"net/http"
	api "./api"
	handler "./handlers"
	helpers "./helpers"
	structs "./structs"
	metrics "./metrics"
	"github.com/swaggo/http-swagger"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/mattn/go-sqlite3"
	_ "./docs"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)
// @title MiniTwit Swagger API
// @version 1.0
// @description Swagger API for Golang Project MiniTwit.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email dallasmaxwell@outlook.com

// @license.name Apache 2.0
// @host localhost:5000

// @BasePath /api/v1

var db *gorm.DB

func init() {
	handler.LoadTemplates()
}

var metricsMonitor = metrics.Combine(
	metrics.HTTPResponseCodeMonitor,
	metrics.HTTPResponseTimeMonitor,
	metrics.HTTPRequestCountMonitor,
)


func main() {

	db := helpers.InitDB()
	defer db.Close()

	db.AutoMigrate(&structs.User{}, &structs.Follower{}, &structs.Message{})

	router := mux.NewRouter()

	router.Handle("/metrics", promhttp.Handler())

	//Får en 404 page not found her
	//Prøver at tage udgangspunkt i https://github.com/swaggo/http-swagger
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

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

	router.HandleFunc("/publictimeline/more", handler.PublicTimelineLoadMore).Methods("GET")

	router.HandleFunc("/{username}", handler.UserpageRoute).Methods("GET")

	router.HandleFunc("/{username}/follow", handler.UserFollowHandler)
	router.HandleFunc("/{username}/unfollow", handler.UserUnfollowHandler)

	apiRoute := mux.NewRouter()
	//apiRoute.HandleFunc("/test", api.Test)
	apiRoute.HandleFunc("/latest", metricsMonitor(api.Get_latest))
	apiRoute.HandleFunc("/register", metricsMonitor(api.Register)).Methods("POST")
	apiRoute.HandleFunc("/msgs", metricsMonitor(api.Messages))
	apiRoute.HandleFunc("/msgs/{username}", metricsMonitor(api.Messages_per_user))
	apiRoute.HandleFunc("/fllws/{username}", metricsMonitor(api.Follow))

	port := 5000
	log.Printf("Server starting on port %v\n", port)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router)) }()

	apiport := 5001
	log.Printf("Api Server starting on port %v\n", apiport)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", apiport), apiRoute)) }()

	go metrics.HTTPRequestCounter()

	select {}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/public/favicon.ico")
}

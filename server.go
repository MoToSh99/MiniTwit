package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	api "./api"
	handler "./handlers"
	helpers "./helpers"
	structs "./structs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
)

var db *gorm.DB

func init() {
	handler.LoadTemplates()
}

var (
	CPU_GAUGE  = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "minitwit_cpu_load_percent",
			Help: "Current load of the CPU in percent.",
	})
	
	REPONSE_COUNTER = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The count of HTTP responses sent.",
	})

	REQ_DURATION_SUMMARY = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "minitwit_request_duration_milliseconds",
		Help: "Request duration distribution.",
	})
)

func main() {
	rand.Seed(time.Now().Unix())

	http.Handle("/metrics", promhttp.Handler())



	go func() {
		for {
				CPU_GAUGE.Add(rand.Float64()* 15 - 5)
				REPONSE_COUNTER.Add(rand.Float64()* 5)
				REQ_DURATION_SUMMARY.Observe(rand.Float64() * 10)

				time.Sleep(time.Second)
		}
}()
	

	db := helpers.InitDB()
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

	router.HandleFunc("/publictimeline/more", handler.PublicTimelineLoadMore).Methods("GET")

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
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router)) }()

	apiport := 5001
	log.Printf("Api Server starting on port %v\n", apiport)
	go func() { log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", apiport), apiRoute)) }()

	select {}
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/public/favicon.ico")
}

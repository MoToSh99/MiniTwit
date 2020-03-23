package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/gorilla/mux"

	structs "./structs"
	handler "./handlers"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

// Helper functions

func getServer() (s *Server) {
	// Before each test, set up a blank database
	os.Remove("/tmp/test.db")
	db, _ := gorm.Open("sqlite3", "/tmp/test.db")

	db.AutoMigrate(&structs.User{}, &structs.Follower{}, &structs.Message{})
	router := mux.NewRouter()
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

	s = &Server{DB: db, Router: router}
	
	return s
}

func getHTMLTemplate(t *testing.T, resp httptest.ResponseRecorder) string {
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	HTML := string(bodyBytes)
	return HTML
}

func register(username string, password string, password2 string, email string, server *Server) httptest.ResponseRecorder {
	// Helper function to register a user
	form := url.Values{}
	if password2 == "" {
		password2 = password
	}
	if email == "" {
		email = username + "@example.com"
	}
	request, _ := http.NewRequest("POST", "/register?username="+username+"&email="+email+"&password="+password+"&password2="+password2, strings.NewReader(form.Encode()))
	response := httptest.NewRecorder()
	server.Router.ServeHTTP(response, request)
	return *response
}

func login(username string, password string, server *Server) httptest.ResponseRecorder {
	// Helper function to login
	form := url.Values{}
	request, _ := http.NewRequest("POST", "/signin?username="+username+"&password="+password, strings.NewReader(form.Encode()))
	response := httptest.NewRecorder()
	server.Router.ServeHTTP(response, request)
	return *response
}

func registerAndLogin(username string, password string, password2 string, email string, server *Server) httptest.ResponseRecorder {
	// Registers and logs in in one go
	form := url.Values{}
	register(username, password, password2, email, server)
	request, _ := http.NewRequest("POST", "/signin?username="+username+"&password="+password, strings.NewReader(form.Encode()))
	response := httptest.NewRecorder()
	server.Router.ServeHTTP(response, request)
	return *response
}

func logout(server *Server) httptest.ResponseRecorder {
	// Helper function to logout
	request, _ := http.NewRequest("GET", "/signout", nil)
	response := httptest.NewRecorder()
	server.Router.ServeHTTP(response, request)
	return *response
}

func addMessage(text string, server *Server) httptest.ResponseRecorder {
	// Records a message
	request, _ := http.NewRequest("POST", "/personaltimeline?text="+text, nil)
	response := httptest.NewRecorder()
	server.Router.ServeHTTP(response, request)
	return *response
}

// Testing functions

func TestRegister(t *testing.T) {
	// Make sure registering works
	server := getServer()

	response := register("foo", "pass1", "pass2", "email", server)
	assert.Equal(t, 200, response.Code, "Ok response is expected")

	response = register("foo", "pass1", "pass2", "email", server)
	html := getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, ("You have to enter a username")))

	response = register("foo", "", "", "email", server)
	html = getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, ("You have to enter a password")))

	response = register("foo", "aa", "bb", "email", server)
	html = getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, ("The two passwords do not match")))

	response = register("foo", "aa", "aa", "", server)
	html = getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, ("You have to enter a valid email address")))

	defer server.DB.Close()
}

func TestLoginLogout(t *testing.T) {
	// Make sure logging in and logging out works
	server := getServer()

	response := registerAndLogin("foo", "default", "default", "example@hotmail.com", server)
	assert.Equal(t, 302, response.Code, "Status found")

	response2 := logout(server)
	assert.Equal(t, 302, response2.Code, "Status found")

	response3 := login("foo", "wrongpassword", server)
	html := getHTMLTemplate(t, response3)
	assert.True(t, true, strings.Contains(html, ("Invalid password")))

	response4 := login("bar", "wrongpassword", server)
	html = getHTMLTemplate(t, response4)
	assert.True(t, true, strings.Contains(html, ("Invalid username")))

	defer server.DB.Close()
}

func TestMessageRecording(t *testing.T) {
	// Check if adding messages works
	server := getServer()

	response := registerAndLogin("foo", "default", "default", "example@hotmail.com", server)
	assert.Equal(t, 302, response.Code, "Status found")

	addMessage("foo bar 123", server)
	addMessage("hello world 123", server)

	request, _ := http.NewRequest("GET", "/public", nil)
	response = *httptest.NewRecorder()
	server.Router.ServeHTTP(&response, request)
	html := getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, "foo bar 123"))
	assert.True(t, true, strings.Contains(html, "hello world 123"))

	defer server.DB.Close()
}

func TestTimelines(t *testing.T) {
	// Make sure that timelines work
	server := getServer()

	response := registerAndLogin("foo", "default", "default", "example@hotmail.com", server)
	assert.Equal(t, 302, response.Code, "Status found")

	addMessage("the message by foo", server)
	request, _ := http.NewRequest("GET", "/public", nil)
	response = *httptest.NewRecorder()
	server.Router.ServeHTTP(&response, request)
	html := getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, "the message by foo"))

	response = logout(server)
	assert.Equal(t, 302, response.Code, "Status found")

	response = registerAndLogin("bar", "default", "default", "example@hotmail.com", server)
	assert.Equal(t, 302, response.Code, "Status found")

	addMessage("the message by bar", server)
	request, _ = http.NewRequest("GET", "/public", nil)
	response = *httptest.NewRecorder()
	server.Router.ServeHTTP(&response, request)
	html = getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, "the message by bar"))

	request, _ = http.NewRequest("GET", "/foo/follow", nil)
	response = *httptest.NewRecorder()
	server.Router.ServeHTTP(&response, request)
	html = getHTMLTemplate(t, response)
	assert.True(t, true, strings.Contains(html, "You are currently following this user."))
	assert.Equal(t, 302, response.Code, "Status found")

	request, _ = http.NewRequest("GET", "/foo/unfollow", nil)
	response = *httptest.NewRecorder()
	server.Router.ServeHTTP(&response, request)
	assert.Equal(t, 302, response.Code, "Status found")

	defer server.DB.Close()
}
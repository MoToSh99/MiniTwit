package cookies

import (
	"net/http"
	"github.com/gorilla/securecookie"
)

var CookieHandler = securecookie.New(
    securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func SetCookie(userName string, res http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := CookieHandler.Encode("cookie", value); err == nil {
		cookie := &http.Cookie{
			Name:  "cookie",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(res, cookie)
	}
}
	
func ClearCookie(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "cookie",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}

package api

import(
    "fmt"
    "net/http"
)
 
func Register(res http.ResponseWriter, req *http.Request) {
    fmt.Fprintln(res, "You have arrived at the API Server!")
}
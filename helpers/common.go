package helpers

import (
	"time"
)

func IsEmpty(data string) bool {
    if len(data) <= 0 {
        return true
    } else {
        return false
    }
}

func GetCurrentTime() string{
	dt := time.Now()
	return (dt.Format("15:04:05 02-01-2006"))
}

func GetConnString() string {
    var db *sql.DB
    var server = "minitwitserver.database.windows.net"
    var port = 1433
    var user = "Minitwit"
    var password = "ITU2020!"
    var database = "minitwitdb"
    
    var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
            server, user, password, port, database)
    return connString
}
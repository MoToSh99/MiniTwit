package main

import (
	"fmt"
	"os"
	"log"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var helpString = "ITU-Minitwit Tweet Flagging ToolUsage:\n\n" +
	"  flag_tool <tweet_id>...\n" +
	"  flag_tool -i\n" +
	"  flag_tool -h\n" +
	"Options:\n" +
	"-h            Show this screen.\n" +
	"-i            Dump all tweets and authors to STDOUT.\n"

func main() {
	db, err := sql.Open("sqlite3", "/tmp/minitwit.db")
	if err != nil {
		log.Fatal(err)
	}

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "-h" {
		fmt.Println(helpString)
	}
	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "-i" {
		query := "SELECT * FROM message"
		/* Execute SQL statement */
		/*rc = sqlite3_exec(db, query, callback, (void *)data, &zErrMsg);
		    if (rc != SQLITE_OK) {
		      fprintf(stderr, "SQL error: %s\n", zErrMsg);
			  sqlite3_free(zErrMsg);*/
		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(rows)
	}
}

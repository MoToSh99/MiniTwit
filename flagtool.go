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
	defer db.Close()

	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "-h" {
		fmt.Println(helpString)
	}

	if len(argsWithoutProg) == 0 {
		fmt.Println(helpString)
	}

	if len(argsWithoutProg) == 1 && argsWithoutProg[0] == "-i" {
		query := "SELECT * FROM message"
		rows, err := db.Query(query)

		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next(){
			//These names are dependant on column names, so ignore _ warning
			var message_id int
			var author_id int
			var text string
			var pub_date string
			var flagged int

			err = rows.Scan(&message_id, &author_id, &text, &pub_date, &flagged)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(message_id,author_id,text,pub_date,flagged)
		}
		
	}

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] != "-i"  && argsWithoutProg[0] != "-h" {
		for i := 1; i < len(argsWithoutProg)+1; i++ {
			query := "Update message Set flagged=1 Where message_id="+argsWithoutProg[i-1]
			_,err := db.Exec(query)
			if err!= nil {
				log.Fatal(err)
				fmt.Println("error with flagging")
			}else{
				fmt.Println("Flagged entry: " + argsWithoutProg[i-1] + "\n")
			}
		}
	}
}

package main

import (
	"fmt"
	"log"

	"github.com/Phillip-England/bible-bot/module/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// operation.Pull("./bible_html", 200)
	// operation.GetCommands("./bible_html", "./bible_commands")
	// operation.MakeBibleJson("./bible_commands", "./bible_json")
	// operation.MakeBibleDb("./bible_json")

	// Open SQLite database
	db, err := sqlx.Open("sqlite3", "./bible.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	verse, err := database.GetVerse(db, "KJV", "JHN", 3, 16)
	if err != nil {
		panic(err)
	}

	fmt.Println(verse)

}

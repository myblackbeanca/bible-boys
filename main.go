package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Phillip-England/bible-bot/module/database"
	"github.com/Phillip-England/bible-bot/module/operation"
	"github.com/Phillip-England/vbf"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	// operation.GetCommands("./bible_html", "./bible_commands")
	// operation.MakeBibleJson("./bible_commands", "./bible_json")
	operation.MakeBibleDb("./bible_json")

	db, err := sqlx.Open("sqlite3", "bible.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux, gCtx := vbf.VeryBestFramework()

	gCtx["DB"] = db

	vbf.AddRoute("GET /{translation}/{book}/{chapter}/{verse}", mux, gCtx, func(w http.ResponseWriter, r *http.Request) {

		db, ok := gCtx["DB"].(*sqlx.DB)
		if !ok {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		trans := r.PathValue("translation")
		book := r.PathValue("book")
		chapterStr := r.PathValue("chapter")
		verseStr := r.PathValue("verse")

		chapterNum, err := strconv.Atoi(chapterStr)
		if err != nil {
			vbf.WriteJSON(w, 400, map[string]interface{}{
				"message": "chapter must be a valid number",
			})
			return
		}

		verseNum, err := strconv.Atoi(verseStr)
		if err != nil {
			vbf.WriteJSON(w, 400, map[string]interface{}{
				"message": "verse must be a valid number",
			})
			return
		}

		verse, err := database.GetVerse(db, trans, book, uint(chapterNum), uint(verseNum))
		if err != nil {
			vbf.WriteJSON(w, 400, map[string]interface{}{
				"message": fmt.Sprintf(`verse number %d does not exist`, verseNum),
			})
			return
		}

		vbf.WriteJSON(w, 200, verse)
	}, vbf.MwCORS, vbf.MwLogger)

	err = vbf.Serve(mux, "8080")
	if err != nil {
		panic(err)
	}

}

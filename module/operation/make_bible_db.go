package operation

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type TranslationDb struct {
	ID           uint   `db:"id"`
	Translation  string `db:"translation"`
	Abbreviation string `db:"abbreviation"`
}

type BookDb struct {
	ID            uint   `db:"id"`
	TranslationID uint   `db:"translation_id"`
	Name          string `db:"name"`
	Abbreviation  string `db:"abbreviation"`
}

type ChapterDb struct {
	ID     uint `db:"id"`
	BookID uint `db:"book_id"`
	Number int  `db:"number"`
}

type VerseDb struct {
	ID        uint   `db:"id"`
	ChapterID uint   `db:"chapter_id"`
	Number    int    `db:"number"`
	Text      string `db:"text"`
	Section   string `db:"section"`
}

const schema = `
CREATE TABLE IF NOT EXISTS translations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  translation TEXT NOT NULL,
  abbreviation TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS books (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  abbreviation TEXT NOT NULL,
  FOREIGN KEY (translation_id) REFERENCES translations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chapters (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  book_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS verses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  chapter_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  text TEXT NOT NULL,
  section TEXT NOT NULL,
  FOREIGN KEY (chapter_id) REFERENCES chapters(id) ON DELETE CASCADE
);`

func MakeBibleDb(bibleJsonSrc string) {

	fmt.Println("Loading db...")
	db, err := sqlx.Open("sqlite3", "./bible.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Creating sql tables...")
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loading in verses...")
	verses := LoadBibleJson(bibleJsonSrc)

	for _, verse := range verses {

		var t TranslationDb
		var b BookDb
		var ch ChapterDb

		fmt.Printf("Checking if the %s translation exists...\n", verse.TranslationAbbv)
		err = db.Get(&t, `SELECT * FROM translations WHERE abbreviation = ?`, verse.TranslationAbbv)
		if err != nil {
			fmt.Printf("The %s translation does not exist, creating it...\n", verse.TranslationAbbv)
			tQuery := `INSERT INTO translations (translation, abbreviation) VALUES (?, ?) RETURNING id`
			err = db.QueryRow(tQuery, verse.Translation, verse.TranslationAbbv).Scan(&t.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created translation %s with ID %d\n", verse.TranslationAbbv, t.ID)
		} else {
			fmt.Printf("Translation %s found with ID %d\n", verse.TranslationAbbv, t.ID)
		}

		fmt.Printf("Checking if book %s exists in translation ID %d...\n", verse.Book, t.ID)
		err = db.Get(&b, `SELECT * FROM books WHERE translation_id = ? AND name = ?`, t.ID, verse.Book)
		if err != nil {
			fmt.Printf("Book %s does not exist, creating it...\n", verse.Book)
			bQuery := `INSERT INTO books (translation_id, name, abbreviation) VALUES (?, ?, ?) RETURNING id`
			err = db.QueryRow(bQuery, t.ID, verse.Book, verse.BookAbbv).Scan(&b.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created book %s with ID %d\n", verse.Book, b.ID)
		} else {
			fmt.Printf("Book %s found with ID %d\n", verse.Book, b.ID)
		}

		fmt.Printf("Checking if chapter %d exists in book ID %d...\n", verse.Chapter, b.ID)
		err = db.Get(&ch, `SELECT * FROM chapters WHERE book_id = ? AND number = ?`, b.ID, verse.Chapter)
		if err != nil {
			fmt.Printf("Chapter %d does not exist, creating it...\n", verse.Chapter)
			chQuery := `INSERT INTO chapters (book_id, number) VALUES (?, ?) RETURNING id`
			err = db.QueryRow(chQuery, b.ID, verse.Chapter).Scan(&ch.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created chapter %d with ID %d\n", verse.Chapter, ch.ID)
		} else {
			fmt.Printf("Chapter %d found with ID %d\n", verse.Chapter, ch.ID)
		}

		fmt.Printf("Inserting verse %d into chapter ID %d...\n", verse.Number, ch.ID)
		vQuery := `INSERT INTO verses (chapter_id, number, text, section) VALUES (?, ?, ?, ?)`
		_, err = db.Exec(vQuery, ch.ID, verse.Number, verse.Text, verse.Section)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted verse %d into chapter ID %d\n", verse.Number, ch.ID)
	}
}

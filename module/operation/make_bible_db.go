package operation

import (
	"fmt"
	"log"

	"github.com/Phillip-England/bible-bot/module/database"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS translation (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  translation TEXT NOT NULL,
  abbreviation TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS book (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  translation_id INTEGER NOT NULL,
  name TEXT NOT NULL,
  abbreviation TEXT NOT NULL,
  FOREIGN KEY (translation_id) REFERENCES translation(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS chapter (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  book_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS verse (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  chapter_id INTEGER NOT NULL,
  number INTEGER NOT NULL,
  text TEXT NOT NULL,
  section TEXT NOT NULL,
  FOREIGN KEY (chapter_id) REFERENCES chapter(id) ON DELETE CASCADE
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

	fmt.Println("Indexing the database for better performance...")
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_book_translation ON book (translation_id);
		CREATE INDEX IF NOT EXISTS idx_book_name ON book (name);
		CREATE INDEX IF NOT EXISTS idx_chapter_book ON chapter (book_id);
		CREATE INDEX IF NOT EXISTS idx_chapter_number ON chapter (book_id, number);
		CREATE INDEX IF NOT EXISTS idx_verse_chapter ON verse (chapter_id);
		CREATE INDEX IF NOT EXISTS idx_verse_number ON verse (chapter_id, number);
`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Loading in verses...")
	verses := LoadBibleJson(bibleJsonSrc)

	for _, verse := range verses {

		var t database.Translation
		var b database.Book
		var ch database.Chapter

		fmt.Printf("Checking if the %s translation exists...\n", verse.TranslationAbbv)
		err = db.Get(&t, `SELECT * FROM translation WHERE abbreviation = ?`, verse.TranslationAbbv)
		if err != nil {
			fmt.Printf("The %s translation does not exist, creating it...\n", verse.TranslationAbbv)
			tQuery := `INSERT INTO translation (translation, abbreviation) VALUES (?, ?) RETURNING id`
			err = db.QueryRow(tQuery, verse.Translation, verse.TranslationAbbv).Scan(&t.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created translation %s with ID %d\n", verse.TranslationAbbv, t.ID)
		} else {
			fmt.Printf("Translation %s found with ID %d\n", verse.TranslationAbbv, t.ID)
		}

		fmt.Printf("Checking if book %s exists in translation ID %d...\n", verse.Book, t.ID)
		err = db.Get(&b, `SELECT * FROM book WHERE translation_id = ? AND name = ?`, t.ID, verse.Book)
		if err != nil {
			fmt.Printf("Book %s does not exist, creating it...\n", verse.Book)
			bQuery := `INSERT INTO book (translation_id, name, abbreviation) VALUES (?, ?, ?) RETURNING id`
			err = db.QueryRow(bQuery, t.ID, verse.Book, verse.BookAbbv).Scan(&b.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created book %s with ID %d\n", verse.Book, b.ID)
		} else {
			fmt.Printf("Book %s found with ID %d\n", verse.Book, b.ID)
		}

		fmt.Printf("Checking if chapter %d exists in book ID %d...\n", verse.Chapter, b.ID)
		err = db.Get(&ch, `SELECT * FROM chapter WHERE book_id = ? AND number = ?`, b.ID, verse.Chapter)
		if err != nil {
			fmt.Printf("Chapter %d does not exist, creating it...\n", verse.Chapter)
			chQuery := `INSERT INTO chapter (book_id, number) VALUES (?, ?) RETURNING id`
			err = db.QueryRow(chQuery, b.ID, verse.Chapter).Scan(&ch.ID)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Created chapter %d with ID %d\n", verse.Chapter, ch.ID)
		} else {
			fmt.Printf("Chapter %d found with ID %d\n", verse.Chapter, ch.ID)
		}

		fmt.Printf("Inserting verse %d into chapter ID %d...\n", verse.Number, ch.ID)
		vQuery := `INSERT INTO verse (chapter_id, number, text, section) VALUES (?, ?, ?, ?)`
		_, err = db.Exec(vQuery, ch.ID, verse.Number, verse.Text, verse.Section)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Inserted verse %d into chapter ID %d\n", verse.Number, ch.ID)
	}
}

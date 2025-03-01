package database

import (
	"github.com/jmoiron/sqlx"
)

func GetAllTranslation(db *sqlx.DB) ([]Translation, error) {
	var translations []Translation
	err := db.Select(&translations, "SELECT id, translation, abbreviation FROM translation")
	if err != nil {
		return translations, err
	}
	return translations, nil
}

func GetTranslationIDByAbbv(db *sqlx.DB, translationAbbv string) (uint, error) {
	var transID uint
	err := db.Get(&transID, "SELECT id FROM translation WHERE abbreviation = ?", translationAbbv)
	if err != nil {
		return transID, err
	}
	return transID, nil
}

func GetBookIDByName(db *sqlx.DB, bookName string) (uint, error) {
	var bookID uint
	err := db.Get(&bookID, "SELECT id FROM book WHERE name = ?", bookName)
	if err != nil {
		return bookID, err
	}
	return bookID, nil
}

func GetBookIDByAbbv(db *sqlx.DB, bookAbbv string) (uint, error) {
	var bookID uint
	err := db.Get(&bookID, "SELECT id FROM book WHERE abbreviation = ?", bookAbbv)
	if err != nil {
		return bookID, err
	}
	return bookID, nil
}

func GetChapterID(db *sqlx.DB, bookID uint, chapterNumber uint) (uint, error) {
	var chapterID uint
	err := db.Get(&chapterID, "SELECT id FROM chapter WHERE book_id = ? AND number = ?", bookID, chapterNumber)
	if err != nil {
		return chapterID, err
	}
	return chapterID, nil
}

func GetVerseByChapterID(db *sqlx.DB, chapterID uint, verseNumber uint) (Verse, error) {
	var verse Verse
	err := db.Get(&verse, "SELECT id, chapter_id, number, text, section FROM verse WHERE chapter_id = ? AND number = ?", chapterID, verseNumber)
	if err != nil {
		return verse, err
	}
	return verse, nil
}

func GetVerse(db *sqlx.DB, translationAbbv string, bookAbbv string, chapterNumber uint, verseNumber uint) (Verse, error) {
	var verse Verse
	query := `
		SELECT v.id, v.chapter_id, v.number, v.text, v.section
		FROM verse v
		JOIN chapter c ON v.chapter_id = c.id
		JOIN book b ON c.book_id = b.id
		JOIN translation t ON b.translation_id = t.id
		WHERE t.abbreviation = ? 
		AND b.abbreviation = ? 
		AND c.number = ? 
		AND v.number = ?;
	`
	err := db.Get(&verse, query, translationAbbv, bookAbbv, chapterNumber, verseNumber)
	if err != nil {
		return verse, err
	}
	return verse, nil
}

func GetChapterWithVerses(db *sqlx.DB, translationAbbv string, bookAbbv string, chapterNumber uint) (Chapter, []Verse, error) {
	var chapter Chapter
	var verses []Verse

	query := `
		SELECT c.id, c.book_id, c.number
		FROM chapter c
		JOIN book b ON c.book_id = b.id
		JOIN translation t ON b.translation_id = t.id
		WHERE t.abbreviation = ? AND b.abbreviation = ? AND c.number = ?
		LIMIT 1;
	`
	err := db.Get(&chapter, query, translationAbbv, bookAbbv, chapterNumber)
	if err != nil {
		return chapter, verses, err
	}

	queryVerses := `
		SELECT v.id, v.chapter_id, v.number, v.text, v.section
		FROM verse v
		WHERE v.chapter_id = ?
		ORDER BY v.number;
	`
	err = db.Select(&verses, queryVerses, chapter.ID)
	if err != nil {
		return chapter, verses, err
	}

	return chapter, verses, nil
}

func GetBookWithChaptersAndVerses(db *sqlx.DB, translationAbbv string, bookAbbv string) (Book, []Chapter, map[uint][]Verse, error) {
	var book Book
	var chapters []Chapter
	versesMap := make(map[uint][]Verse)

	queryBook := `
		SELECT b.id, b.translation_id, b.name, b.abbreviation
		FROM book b
		JOIN translation t ON b.translation_id = t.id
		WHERE t.abbreviation = ? AND b.abbreviation = ?
		LIMIT 1;
	`
	err := db.Get(&book, queryBook, translationAbbv, bookAbbv)
	if err != nil {
		return book, chapters, versesMap, err
	}

	queryChapters := `
		SELECT c.id, c.book_id, c.number
		FROM chapter c
		WHERE c.book_id = ?
		ORDER BY c.number;
	`
	err = db.Select(&chapters, queryChapters, book.ID)
	if err != nil {
		return book, chapters, versesMap, err
	}

	queryVerses := `
		SELECT v.id, v.chapter_id, v.number, v.text, v.section
		FROM verse v
		WHERE v.chapter_id = ?
		ORDER BY v.number;
	`

	for _, chapter := range chapters {
		var verses []Verse
		err = db.Select(&verses, queryVerses, chapter.ID)
		if err != nil {
			return book, chapters, versesMap, err
		}
		versesMap[chapter.ID] = verses
	}

	return book, chapters, versesMap, nil
}

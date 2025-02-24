package global

import "fmt"

const BIBLE_SRC = "https://bible.com/bible"

func GetBibleUrl(book Book, translation Translation, chapter int) string {
	return fmt.Sprintf(`%s/%d/%s.%d.%s`, BIBLE_SRC, translation.PageCode, book.Abbreviation, chapter, translation.Abbreviation)
}

func GetBibleOut(out string, book Book, translation Translation) string {
	return fmt.Sprintf(`%s/%s/%s`, out, translation.Abbreviation, book.Abbreviation)
}
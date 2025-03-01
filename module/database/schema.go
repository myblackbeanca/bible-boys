package database

type Translation struct {
	ID           uint   `db:"id"`
	Translation  string `db:"translation"`
	Abbreviation string `db:"abbreviation"`
}

type Book struct {
	ID            uint   `db:"id"`
	TranslationID uint   `db:"translation_id"`
	Name          string `db:"name"`
	Abbreviation  string `db:"abbreviation"`
}

type Chapter struct {
	ID     uint `db:"id"`
	BookID uint `db:"book_id"`
	Number int  `db:"number"`
}

type Verse struct {
	ID        uint   `db:"id"`
	ChapterID uint   `db:"chapter_id"`
	Number    int    `db:"number"`
	Text      string `db:"text"`
	Section   string `db:"section"`
}

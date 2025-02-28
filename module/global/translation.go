package global

type Translation struct {
	Abbreviation string
	Name         string
	PageCode     int32
}

func GetBibleTranslations() []Translation {
	return []Translation{
		{"KJV", "King James Version", 1},
		{"ASV", "American Standard Version", 12},
		{"WEBUS", "World English Bible", 206},
	}
}

func GetTranslationByAbbreviation(abbr string) (Translation, bool) {
	for _, t := range GetBibleTranslations() {
		if t.Abbreviation == abbr {
			return t, true
		}
	}
	return Translation{}, false
}

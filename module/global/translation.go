package global

type Translation struct {
	Abbreviation string
	Name         string
	PageCode     int32
}

func GetBibleTranslations() []Translation {
	return []Translation{
		{"KJV", "King James Version", 1},
		{"NKJV", "New King James Version", 114},
		{"NIV", "New International Version", 111},
		{"ESV", "English Standard Version", 59},
		{"NLT", "New Living Translation", 116},
		{"CSB", "Christian Standard Bible", 1713},
		{"RSV", "Revised Standard Version", 2020},
		{"NRSV", "New Revised Standard Version", 2016},
		{"ASV", "American Standard Version", 12},
		{"AMP", "Amplified Bible", 1588},
		{"CEB", "Common English Bible", 37},
		{"CEV", "Contemporary English Version", 392},
		{"GNT", "Good News Translation", 68},
		{"HCSB", "Holman Christian Standard Bible", 72},
		{"JUB", "Jubilee Bible", 1077},
		{"MEV", "Modern English Version", 1171},
		{"NCV", "New Century Version", 105},
		{"NABRE", "New American Bible Revised Edition", 436},
		{"NET", "New English Translation", 107},
		{"ERV", "Easy-to-Read Version", 406},
		{"MSG", "The Message", 97},
		{"DARBY", "Darby Translation", 478},
	}
}
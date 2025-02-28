package operation

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
)

func LoadBibleJson(src string) []Verse {
	var verses []Verse
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		verse := &Verse{}
		err = json.Unmarshal(fileBytes, verse)
		if err != nil {
			return err
		}
		verses = append(verses, *verse)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return verses
}

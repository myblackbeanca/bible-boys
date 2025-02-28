package operation

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Phillip-England/bible-bot/module/global"
)

type Verse struct {
	Number          int    `json:"number"`
	Text            string `json:"text"`
	Section         string `json:"section"`
	Translation     string `json:"translation"`
	TranslationAbbv string `json:"translation_abbreviation"`
	Book            string `json:"book"`
	BookAbbv        string `json:"book_abbreviation"`
	Chapter         int    `json:"chapter"`
}

func MakeBibleJson(commandSrc string, out string) {
	err := os.RemoveAll(out)
	if err != nil {
		panic(err)
	}
	var verses []Verse
	err = filepath.Walk(commandSrc, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fileStr := string(fileBytes)
		commands := strings.Split(fileStr, "\n")
		verse := &Verse{}
		currentVerseNumber := ""
		for _, command := range commands {
			parts := strings.Split(command, ":")
			commandName := parts[0]
			commandValue := parts[1]
			if commandName == "TRANSLATION" {
				verse.TranslationAbbv = commandValue
			}
			if commandName == "BOOK" {
				verse.BookAbbv = commandValue
			}
			if commandName == "CHAPTER" {
				chapterNum, err := strconv.Atoi(commandValue)
				if err != nil {
					return err
				}
				verse.Chapter = chapterNum
			}
			if commandName == "SECTION" {
				verse.Section = commandValue
			}
			if commandName == "LABEL" {
				currentVerseNumber = commandValue
			}
			if commandName == "VERSE" {
				verseNum, err := strconv.Atoi(currentVerseNumber)
				if err != nil {
					return err
				}
				verse.Number = verseNum
				verse.Text = commandValue
			}
			if verse.BookAbbv != "" && verse.Chapter != 0 && verse.Number != 0 && verse.Text != "" && verse.TranslationAbbv != "" {
				book, _ := global.GetBookByAbbreviation(verse.BookAbbv)
				verse.Book = book.Name
				trans, _ := global.GetTranslationByAbbreviation(verse.TranslationAbbv)
				verse.Translation = trans.Name
				verses = append(verses, *verse)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, verse := range verses {
		jsonData, err := json.MarshalIndent(verse, "", " ")
		if err != nil {
			panic(err)
		}
		outDir := fmt.Sprintf(`%s/%s/%s/%d`, out, verse.TranslationAbbv, verse.BookAbbv, verse.Chapter)
		out := fmt.Sprintf(`%s/%d.json`, outDir, verse.Number)
		fmt.Println(out)
		err = os.MkdirAll(outDir, 0755)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(out, jsonData, 0644)
		if err != nil {
			panic(err)
		}

	}

}

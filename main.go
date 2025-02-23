package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

const maxConcurrent = 30

func main() {

	err := os.RemoveAll("./bible")
	if err != nil {
		panic(err)
	}

	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup

	for _, b := range getBibleBooks() {
		for _, t := range getBibleTranslations() {
			for chapter := 1; chapter < 60; chapter++ {
				url := getBibleUrl(b, t, chapter)
				out := getBibleOut(b, t)
				outHtml := fmt.Sprintf(`%s/%d.html`, out, chapter)

				wg.Add(1)
				sem <- struct{}{}

				go func(url, out string) {
					defer wg.Done()
					defer func() { <-sem }()

					fmt.Println("requesting: ", url)
					resp, err := http.Get(url)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					err = os.MkdirAll(out, 0755)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}

					fmt.Println("writing to: ", out)
					err = os.WriteFile(outHtml, body, 0644)
					if err != nil {
						fmt.Println("Error:", err)
						return
					}
				}(url, out)
			}
		}
	}

	wg.Wait() // Wait for all goroutines to finish
}

//================================
// const
//================================

const BIBLE_SRC = "https://bible.com/bible"

func getBibleUrl(book Book, translation Translation, chapter int) string {
	return fmt.Sprintf(`%s/%d/%s.%d.%s`, BIBLE_SRC, translation.PageCode, book.Abbreviation, chapter, translation.Abbreviation)
}

func getBibleOut(book Book, translation Translation) string {
	return fmt.Sprintf(`./bible/%s/%s`, translation.Abbreviation, book.Abbreviation)
}

//================================
// bible books
//================================

type Book struct {
	Abbreviation string
	Name         string
}

func getBibleBooks() []Book {
	return []Book{
		{"GEN", "Genesis"},
		{"EXO", "Exodus"},
		{"LEV", "Leviticus"},
		{"NUM", "Numbers"},
		{"DEU", "Deuteronomy"},
		{"JOS", "Joshua"},
		{"JDG", "Judges"},
		{"RUT", "Ruth"},
		{"1SA", "1 Samuel"},
		{"2SA", "2 Samuel"},
		{"1KI", "1 Kings"},
		{"2KI", "2 Kings"},
		{"1CH", "1 Chronicles"},
		{"2CH", "2 Chronicles"},
		{"EZR", "Ezra"},
		{"NEH", "Nehemiah"},
		{"EST", "Esther"},
		{"JOB", "Job"},
		{"PSA", "Psalms"},
		{"PRO", "Proverbs"},
		{"ECC", "Ecclesiastes"},
		{"SNG", "Song of Solomon"},
		{"ISA", "Isaiah"},
		{"JER", "Jeremiah"},
		{"LAM", "Lamentations"},
		{"EZK", "Ezekiel"},
		{"DAN", "Daniel"},
		{"HOS", "Hosea"},
		{"JOL", "Joel"},
		{"AMO", "Amos"},
		{"OBA", "Obadiah"},
		{"JON", "Jonah"},
		{"MIC", "Micah"},
		{"NAM", "Nahum"},
		{"HAB", "Habakkuk"},
		{"ZEP", "Zephaniah"},
		{"HAG", "Haggai"},
		{"ZEC", "Zechariah"},
		{"MAL", "Malachi"},
		{"MAT", "Matthew"},
		{"MRK", "Mark"},
		{"LUK", "Luke"},
		{"JHN", "John"},
		{"ACT", "Acts"},
		{"ROM", "Romans"},
		{"1CO", "1 Corinthians"},
		{"2CO", "2 Corinthians"},
		{"GAL", "Galatians"},
		{"EPH", "Ephesians"},
		{"PHP", "Philippians"},
		{"COL", "Colossians"},
		{"1TH", "1 Thessalonians"},
		{"2TH", "2 Thessalonians"},
		{"1TI", "1 Timothy"},
		{"2TI", "2 Timothy"},
		{"TIT", "Titus"},
		{"PHM", "Philemon"},
		{"HEB", "Hebrews"},
		{"JAS", "James"},
		{"1PE", "1 Peter"},
		{"2PE", "2 Peter"},
		{"1JN", "1 John"},
		{"2JN", "2 John"},
		{"3JN", "3 John"},
		{"JUD", "Jude"},
		{"REV", "Revelation"},
	}
}

//================================
// bible translations
//================================

type Translation struct {
	Abbreviation string
	Name         string
	PageCode     int32
}

func getBibleTranslations() []Translation {
	return []Translation{
		{"KJV", "King James Version", 1},
		{"NKJV", "New King James Version", 114},
		{"NIV", "New International Version", 111},
		{"ESV", "English Standard Version", 59},
		{"NLT", "New Living Translation", 13},
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

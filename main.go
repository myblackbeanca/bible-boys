package main

import (
	"github.com/Phillip-England/bible-bot/module/operation"
)

func main() {
	// operation.Pull("./bible_html", 200)
	// operation.GetCommands("./bible_html", "./bible_commands")
	// operation.MakeBibleJson("./bible_commands", "./bible_json")
	operation.MakeBibleDb("./bible_json")
}

// func testBible() {
// 	verses := operation.LoadBibleJson("./bible_json")
// 	var kjv []operation.Verse
// 	var asv []operation.Verse
// 	var webus []operation.Verse
// 	for _, verse := range verses {
// 		if verse.TranslationAbbv == "KJV" {
// 			kjv = append(kjv, verse)
// 		}
// 		if verse.TranslationAbbv == "ASV" {
// 			asv = append(asv, verse)
// 		}
// 		if verse.TranslationAbbv == "WEBUS" {
// 			webus = append(webus, verse)
// 		}
// 	}
// 	checkDuplicates(kjv)
// 	checkDuplicates(asv)
// 	checkDuplicates(webus)
// 	checkMissingVerses(kjv)
// 	checkMissingVerses(asv)
// 	checkMissingVerses(webus)
// }

// func checkDuplicates(verses []operation.Verse) {
// 	verseMap := make(map[string]int)
// 	duplicates := 0
// 	for _, verse := range verses {
// 		key := fmt.Sprintf("%s-%d-%d-%s", verse.Book, verse.Chapter, verse.Number, verse.Text)
// 		verseMap[key]++
// 		if verseMap[key] > 1 {
// 			fmt.Printf("Duplicate found: %s %d:%d\n", verse.Book, verse.Chapter, verse.Number)
// 			duplicates++
// 		}
// 	}
// 	if duplicates == 0 {
// 		fmt.Println("No duplicate verses found!")
// 	} else {
// 		fmt.Printf("Total duplicates found: %d\n", duplicates)
// 	}
// }

// func checkMissingVerses(verses []operation.Verse) {
// 	chapterVerses := make(map[string]map[int]bool)
// 	for _, verse := range verses {
// 		key := fmt.Sprintf("%s-%d", verse.Book, verse.Chapter)
// 		if chapterVerses[key] == nil {
// 			chapterVerses[key] = make(map[int]bool)
// 		}
// 		chapterVerses[key][verse.Number] = true
// 	}
// 	missingCount := 0
// 	for key, versesMap := range chapterVerses {
// 		var missing []int
// 		maxVerseNum := 0
// 		for verseNum := range versesMap {
// 			if verseNum > maxVerseNum {
// 				maxVerseNum = verseNum
// 			}
// 		}
// 		for i := 1; i <= maxVerseNum; i++ {
// 			if !versesMap[i] {
// 				missing = append(missing, i)
// 			}
// 		}
// 		if len(missing) > 0 {
// 			fmt.Printf("Missing verses in %s: %v\n", key, missing)
// 			missingCount += len(missing)
// 		}
// 	}
// 	if missingCount == 0 {
// 		fmt.Println("No missing verses found!")
// 	} else {
// 		fmt.Printf("Total missing verses: %d\n", missingCount)
// 	}
// }

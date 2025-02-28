package operation

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetCommands(bibleSrc string, out string) {
	if err := os.RemoveAll(out); err != nil {
		panic(err)
	}
	filepath.Walk(bibleSrc, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		pathParts := strings.Split(path, "/")
		t := pathParts[1]
		b := pathParts[2]
		v := strings.Replace(pathParts[3], ".html", "", 1)
		fmt.Println(t, b, v)
		fileHtml, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(fileHtml)))
		if err != nil {
			panic(err)
		}
		doc.Find(".ChapterContent_note__YlDW0").Each(func(i int, s *goquery.Selection) {
			s.Remove()
		})
		commands := []string{
			"TRANSLATION:" + t,
			"BOOK:" + b,
			"CHAPTER:" + v,
		}
		doc.Find("*").Each(func(i int, s *goquery.Selection) {

			className, _ := s.Attr("class")
			text := strings.TrimSpace(s.Text())
			dataUsfm, usfmExists := s.Attr("data-usfm")
			if text == "" {
				return
			}
			if goquery.NodeName(s) == "h1" {
				commands = append(commands, "TITLE:"+text)
			}

			if className == "ChapterContent_heading__xBDcs" {
				commands = append(commands, "SECTION:"+text)
			}

			if usfmExists {
				parts := strings.Split(dataUsfm, ".")
				if len(parts) == 2 {
					return
				}
				verseNumber := parts[len(parts)-1]
				label := "LABEL:" + verseNumber
				foundLabel := false
				for _, value := range commands {
					if value == label {
						foundLabel = true
						break
					}
				}
				if !foundLabel {
					commands = append(commands, label)
				}
				text = strings.Replace(text, verseNumber, "", 1)
				commands = append(commands, "VERSE:"+text)
			}

		})
		bibleSrc = strings.ReplaceAll(bibleSrc, "./", "")
		outPath := strings.ReplaceAll(path, bibleSrc, out)
		outPath = strings.ReplaceAll(outPath, ".html", ".txt")
		parts := strings.Split(outPath, "/")
		dir := strings.Join(parts[:len(parts)-1], "/")
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
		fmt.Println("writing commands to disk at: " + outPath)
		err = os.WriteFile(outPath, []byte(strings.Join(commands, "\n")), 0755)
		if err != nil {
			panic(err)
		}
		return nil
	})
}

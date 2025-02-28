package operation

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/Phillip-England/bible-bot/module/global"
	"github.com/PuerkitoBio/goquery"
)

func Pull(out string, maxConcurrent int) {
	if err := os.RemoveAll(out); err != nil {
		panic(err)
	}
	sem := make(chan struct{}, maxConcurrent)
	var wg sync.WaitGroup
	for _, b := range global.GetBibleBooks() {
		for _, t := range global.GetBibleTranslations() {
			for chapter := 1; chapter < 160; chapter++ {
				book := b
				translation := t
				chap := chapter
				url := global.GetBibleUrl(book, translation, chap)
				dir := global.GetBibleOut(out, book, translation)
				outHtml := fmt.Sprintf("%s/%d.html", dir, chap)
				wg.Add(1)
				sem <- struct{}{}
				go func(url, dir, outHtml string) {
					defer wg.Done()
					defer func() { <-sem }()
					fmt.Println("requesting:", url)
					resp, err := http.Get(url)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						panic(err)
					}
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
					if err != nil {
						panic(err)
					}
					notFoundSel := doc.Find(".ChapterContent_not-avaliable-span__WrOM_")
					if notFoundSel.Length() > 0 {
						return
					}
					contentSel := doc.Find(".ChapterContent_yv-bible-text__tqVMm")
					bodyHtml, err := contentSel.Html()
					if err != nil {
						panic(err)
					}
					if err := os.MkdirAll(dir, 0755); err != nil {
						panic(err)
					}
					fmt.Println("writing to:", outHtml)
					if err := os.WriteFile(outHtml, []byte(bodyHtml), 0644); err != nil {
						panic(err)
					}
				}(url, dir, outHtml)

			}
		}
	}
	wg.Wait()
}

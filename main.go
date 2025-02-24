package main

import (
	"github.com/Phillip-England/bible-bot/module/operation"
)

const maxConcurrent = 300

func main() {
	operation.Pull("./bible_html", 100)
}





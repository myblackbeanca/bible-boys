const jsdom = require("jsdom");
const { JSDOM } = jsdom;

//============================
// consts
//============================

const BIBLE_SRC = 'https://bible.com/bible'

enum BookName {
  GEN = "GEN"
}


enum Translation {
  KJV = 'KJV'
}

enum HTMLClass {
  BOOK = '.ChapterContent_book__VkdB2',
  VERSE = '.ChapterContent_verse__57FIw',
  VERSE_NUMBER = '.ChapterContent_label__R2PLt',
  VERSE_CONTENT = '.ChapterContent_content__RrUqA'
}

//=============================
// utils
//=============================


//=============================
// bible
//=============================

class Bible {
  books: Book[] = []
  constructor() {
    for (const key of Object.keys(BookName)) {
      let book = new Book(key)
      this.books.push(book)
    }
  }
}

//=============================
// books
//=============================

class Book {
  name: string
  constructor(name: string) {
    this.name = name
  }
}

//=============================
// chapters
//=============================





let bible = new Bible()
console.log(bible.books)

// class BibleBot {
//   root: string = "https://bible.com";

//   async scan(book: string, chapter: string, translation: string) {
//     // get the html
//     let endpoint = `${this.root}/bible/1/${book}.${chapter}.${translation}`;
//     let res = await fetch(endpoint);
//     let html = await res.text();

//     // make a doc
//     let dom = new JSDOM(html);
//     let doc = dom.window.document;
//     let verseContent = ''

//     // get the chapter element
//     let chapElm = doc.querySelector(HTMLClass.BOOK);

//     // get the verse elements
//     let verseElms = chapElm.querySelectorAll(HTMLClass.VERSE);

//     // go through each verse elm
//     for (let i = 0; i < verseElms.length; i++) {
//       // get the verse id
//       let verseElm = verseElms[i];
//       let verseID = verseElm.getAttribute("data-usfm");

//       // get the label element (contains the verse number)
//       let labelElm = verseElm.querySelector(HTMLClass.VERSE_NUMBER);
//       if (!labelElm) {
//         continue
//       }

//       // get the verse number
//       let verseNumber = labelElm.textContent
      
//       // get the actual verse from the content element
//       let contentElms = verseElm.querySelectorAll(HTMLClass.VERSE_CONTENT)
      
//       // collect the verse content
//       for (let i2 = 0; i2 < contentElms.length; i2++) {
//         let contentElm = contentElms[i2]
//         let verseText = contentElm.textContent
//         verseContent += verseText + ' '
//       }

//       // trim it up each cycle
//       verseContent.trim()

//     }

//     let file = Bun.file('./job_2.html')
//     await file.write(verseContent)
//   }
// }

// let bot = new BibleBot();
// bot.scan("GEN", "1", "KVJ");

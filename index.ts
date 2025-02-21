const jsdom = require("jsdom");
const { JSDOM } = jsdom;

//============================
// consts
//============================

const BIBLE_SRC = 'https://bible.com/bible'

function BIBLE_URL(bookName: string, chapter: number, translation: string) {
  return BIBLE_SRC+`/1/${bookName}.${chapter}.${translation}`
}

enum BookName {
  GEN = "GEN",
  EXO = "EXO",
  LEV = "LEV",
  NUM = "NUM",
  DEU = "DEU",
  JOS = "JOS",
  JDG = "JDG",
  RUT = "RUT",
  "1SA" = "1SA",
  "2SA" = "2SA",
  "1KI" = "1KI",
  "2KI" = "2KI",
  "1CH" = "1CH",
  "2CH" = "2CH",
  EZR = "EZR",
  NEH = "NEH",
  EST = "EST",
  JOB = "JOB",
  PSA = "PSA",
  PRO = "PRO",
  ECC = "ECC",
  SNG = "SNG",
  ISA = "ISA",
  JER = "JER",
  LAM = "LAM",
  EZK = "EZK",
  DAN = "DAN",
  HOS = "HOS",
  JOL = "JOL",
  AMO = "AMO",
  OBA = "OBA",
  JON = "JON",
  MIC = "MIC",
  NAM = "NAM",
  HAB = "HAB",
  ZEP = "ZEP",
  HAG = "HAG",
  ZEC = "ZEC",
  MAL = "MAL",

  MAT = "MAT",
  MRK = "MRK",
  LUK = "LUK",
  JHN = "JHN",
  ACT = "ACT",
  ROM = "ROM",
  "1CO" = "1CO",
  "2CO" = "2CO",
  GAL = "GAL",
  EPH = "EPH",
  PHP = "PHP",
  COL = "COL",
  "1TH" = "1TH",
  "2TH" = "2TH",
  "1TI" = "1TI",
  "2TI" = "2TI",
  TIT = "TIT",
  PHM = "PHM",
  HEB = "HEB",
  JAS = "JAS",
  "1PE" = "1PE",
  "2PE" = "2PE",
  "1JN" = "1JN",
  "2JN" = "2JN",
  "3JN" = "3JN",
  JUD = "JUD",
  REV = "REV"
}



enum Translation {
  KJV = 'KJV',
  AMP = 'AMP'
}

enum HTMLClass {
  BOOK = 'ChapterContent_book__VkdB2',
  VERSE = 'ChapterContent_verse__57FIw',
  VERSE_NUMBER = 'ChapterContent_label__R2PLt',
  VERSE_CONTENT = 'ChapterContent_content__RrUqA',
  CHAPTER_NOT_FOUND = 'ChapterContent_not-avaliable-span__WrOM_'
}

//=============================
// utils
//=============================


//=============================
// bible
//=============================

class Bible {
  books: Record<string, Book> = {};
  translation: Translation

  private constructor(translation: Translation, books: Record<string, Book>) {
    this.translation = translation
    this.books = books
  }

  static async create(translation: Translation): Promise<Bible> {
    const bookNames = Object.keys(BookName);
    let books: Record<string, Book> = {}
    await Promise.all(bookNames.map(async name => {
      let book = await Book.create(name, translation)
      books[name] = book      
    }));
    return new Bible(translation, books);
  }

  getBook(name: BookName) {
    return this.books[name]
  }
}

//=============================
// books
//=============================

class Book {
  name: string;
  numberOfChapters: number = 0;
  translation: Translation
  requestUrls: string[] = [];
  verses: Verse[] = []

  private constructor(name: string, translation: Translation) {
    this.name = name;
    this.translation = translation
  }

  static async create(name: string, translation: Translation): Promise<Book> {
    const book = new Book(name, translation);
    await book.init();
    return book;
  }

  private async init() {
    await this.initNumberOfChapters();
    this.initRequestUrls();
  }

  // gathers all the urls needed to get the book's verses
  private initRequestUrls() {
    for (let i = 1; i <= this.numberOfChapters; i++) {
      this.requestUrls.push(BIBLE_URL(this.name, i, this.translation));
    }
  }

  // figures out how many chapters are in the book
  private async initNumberOfChapters() {
    let reqUrls = [];
    for (let i = 1; i < 70; i++) {
      reqUrls.push(BIBLE_URL(this.name, i, this.translation));
    }

    await Promise.all(
      reqUrls.map(async (url) => {
        let res = await fetch(url);
        let text = await res.text();
        if (text.includes(HTMLClass.CHAPTER_NOT_FOUND)) {
          return;
        }
        let parts = url.split("/");
        let lastPart = parts[parts.length - 1];
        let moreParts = lastPart.split(".");
        let chapter = moreParts[1];

        if (Number(chapter) > this.numberOfChapters) {
          this.numberOfChapters = Number(chapter);
        }
      })
    );
  }
}

//=============================
// verse
//=============================

class Verse {
  number: number
  constructor(number: number) {
    this.number = number
  }
}

//=============================
// verse
//=============================





(async () => {


  const bible = await Bible.create(Translation.AMP);

  let genesis = bible.getBook(BookName.GEN)

  console.log(genesis)




})();
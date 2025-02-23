const jsdom = require("jsdom");
const { JSDOM } = jsdom;

const BIBLE_SRC = 'https://bible.com/bible'
const CHAPTER_LIMIT = 60

function BIBLE_URL(bookName: string, chapter: number, translation: string) {
  return BIBLE_SRC+`/1/${bookName}.${chapter}.${translation}`
}

const bookNames: Record<string, string> = {
  GEN: "Genesis",
  EXO: "Exodus",
  LEV: "Leviticus",
  NUM: "Numbers",
  DEU: "Deuteronomy",
  JOS: "Joshua",
  JDG: "Judges",
  RUT: "Ruth",
  "1SA": "1 Samuel",
  "2SA": "2 Samuel",
  "1KI": "1 Kings",
  "2KI": "2 Kings",
  "1CH": "1 Chronicles",
  "2CH": "2 Chronicles",
  EZR: "Ezra",
  NEH: "Nehemiah",
  EST: "Esther",
  JOB: "Job",
  PSA: "Psalms",
  PRO: "Proverbs",
  ECC: "Ecclesiastes",
  SNG: "Song of Solomon",
  ISA: "Isaiah",
  JER: "Jeremiah",
  LAM: "Lamentations",
  EZK: "Ezekiel",
  DAN: "Daniel",
  HOS: "Hosea",
  JOL: "Joel",
  AMO: "Amos",
  OBA: "Obadiah",
  JON: "Jonah",
  MIC: "Micah",
  NAM: "Nahum",
  HAB: "Habakkuk",
  ZEP: "Zephaniah",
  HAG: "Haggai",
  ZEC: "Zechariah",
  MAL: "Malachi",
  MAT: "Matthew",
  MRK: "Mark",
  LUK: "Luke",
  JHN: "John",
  ACT: "Acts",
  ROM: "Romans",
  "1CO": "1 Corinthians",
  "2CO": "2 Corinthians",
  GAL: "Galatians",
  EPH: "Ephesians",
  PHP: "Philippians",
  COL: "Colossians",
  "1TH": "1 Thessalonians",
  "2TH": "2 Thessalonians",
  "1TI": "1 Timothy",
  "2TI": "2 Timothy",
  TIT: "Titus",
  PHM: "Philemon",
  HEB: "Hebrews",
  JAS: "James",
  "1PE": "1 Peter",
  "2PE": "2 Peter",
  "1JN": "1 John",
  "2JN": "2 John",
  "3JN": "3 John",
  JUD: "Jude",
  REV: "Revelation"
};


const translations: Record<string, string> = {
  KJV: "King James Version",
  AMP: "Amplified Bible",
  NIV: "New International Version",
  ESV: "English Standard Version",
  NLT: "New Living Translation",
  NKJV: "New King James Version",
  NASB: "New American Standard Bible",
  CSB: "Christian Standard Bible",
  RSV: "Revised Standard Version",
  NRSV: "New Revised Standard Version",
  HCSB: "Holman Christian Standard Bible",
  ASV: "American Standard Version",
  DBY: "Darby Bible",
  YLT: "Young's Literal Translation",
  WEB: "World English Bible",
  CEV: "Contemporary English Version",
  GNT: "Good News Translation",
  MSG: "The Message",
  TPT: "The Passion Translation",
  MEV: "Modern English Version",
  NCV: "New Century Version",
  ERV: "Easy-to-Read Version"
};

let chapterRequests: BibleChapterRequest[] = []

class BibleChapterRequest {
  src: string
  bookKey: string
  bookValue: string
  chapterNumber: number
  translationKey: string
  translationValue: string
  out: string
  constructor(bookKey: string, bookValue: string, chapterNumber: number, translationKey: string, translationValue: string) {
    this.bookKey = bookKey
    this.bookValue = bookValue
    this.chapterNumber = chapterNumber
    this.translationKey = translationKey
    this.translationValue = translationValue
    this.src = `https://bible.com/bible/1/${bookKey}.${chapterNumber}.${translationKey}`
    this.out = `./bible/${translationKey}/${bookKey}/${chapterNumber}.html`
  }
}

Object.entries(bookNames).forEach(([bookKey, bookValue]) => {
  Object.entries(translations).forEach(([translationKey, translationsValue]) => {
    for (let chapter = 1; chapter < CHAPTER_LIMIT; chapter++) {
      let chapterRequest = new BibleChapterRequest(bookKey, bookValue, chapter, translationKey, translationsValue)
      chapterRequests.push(chapterRequest)
    }
  })
})

for (let i = 0; i < chapterRequests.length; i++) {
  let chapter = chapterRequests[i]
  let res = await fetch(chapter.src)
  let html = await res.text()
  let jsdom = new JSDOM(html)
  let document = jsdom.window.document
  let notFoundElement = document.querySelector('.ChapterContent_not-avaliable-span__WrOM_')
  if (notFoundElement) {
    continue
  }
  let file = Bun.file(chapter.out)
  console.log(`writing ${chapter.src}`)
  await file.write(html)
}
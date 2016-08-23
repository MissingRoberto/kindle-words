package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
   "github.com/codegangsta/cli"
   "os"
   "fmt"
  "github.com/jszroberto/kindle-words/kindledb"
)

func main() {

  app := cli.NewApp()
  app.EnableBashCompletion = true
  app.Name = "kindle-words"
  app.Usage = "Provides methods to work with vocabulary builder "

  // global level flags
  app.Flags = []cli.Flag{
    cli.BoolFlag{
      Name:  "verbose",
      Usage: "Show more output",
    },
  }

  // Commands
  app.Commands = []cli.Command{
    {
    Name: "export",
    Flags: []cli.Flag{
      cli.BoolFlag{
        Name:  "html",
        Usage: "exports to html",
      },
      cli.BoolFlag{
        Name:  "vocab.com",
        Usage: "exports to vocab.com",
      },
      cli.BoolFlag{
        Name:  "csv",
        Usage: "export to csv",
      },
      cli.BoolFlag{
        Name:  "evernote",
        Usage: "export to evernote",
      },
    },
    Usage:  "Export words",
    Action: export,
    },
    {
    Name: "delete-book",
    Usage:  "Delete book",
    Action: deleteBook,
    },
    {
    Name: "delete-word",
    Usage:  "Delete word",
    Action: deleteWord,
    },
    {
    Name: "books",
    Usage:  "show book title",
    Action: listBooks,
    },
  }

  app.Run(os.Args)

}

func listBooks(c *cli.Context){
  var language string
  if c.NArg() > 0 {
   language = c.Args()[0]
  }

  kindle, err:= kindledb.NewKindleDB()
  if err != nil {
    log.Fatal(err)
  }
  defer kindle.Close()

  books,err := kindle.ReadBooksInfo()
  if err != nil {
    log.Fatal(err)
  }

  for _,book:= range books{
    if language == "" || book.Language == language{
      fmt.Println(book.Title)
    }
  }
}

func deleteBook(c *cli.Context){
  var bookName string
  if c.NArg() > 0 {
   bookName = c.Args()[0]
  }else{
    cli.ShowSubcommandHelp(c)
    return
  }

  kindle, err:= kindledb.NewKindleDB()
  if err != nil {
    log.Fatal(err)
  }
  defer kindle.Close()

  err = kindle.RemoveBook(bookName)
  if err != nil {
    log.Fatal(err)
  }
}

func deleteWord(c *cli.Context){
  var word string
  if c.NArg() > 0 {
   word = c.Args()[0]
  }else{
    cli.ShowSubcommandHelp(c)
    return
  }

  kindle, err:= kindledb.NewKindleDB()
  if err != nil {
    log.Fatal(err)
  }
  defer kindle.Close()

  err = kindle.RemoveWord(word)
  if err != nil {
    log.Fatal(err)
  }
}

func export(c *cli.Context){

  kindle, err:= kindledb.NewKindleDB()
  if err != nil {
    log.Fatal(err)
  }
  defer kindle.Close()

  words, err := kindle.ReadWords()
  if err != nil {
    log.Fatal(err)
  }

  // if c.NArg() > 0 {
  //  bookName = c.Args()[0]
  //  fmt.Println(bookName)
  // }

  if c.Bool("html") {
    err = exportHtml("./export/html",words)
    if err != nil {
  		log.Fatal(err)
  	}
  }else if c.Bool("vocab.com"){
    err = exportVocabularyCom("./export/vocab.com",words)
    if err != nil {
      log.Fatal(err)
    }
  }else if c.Bool("evernote"){
    err = exportEvernote("./export/evernote",words)
    if err != nil {
      log.Fatal(err)
    }
  }else if c.Bool("csv"){
    err = exportToCSV("./export/result.csv",words)
    if err != nil {
  		log.Fatal(err)
  	}
  }else{
    cli.ShowSubcommandHelp(c)
  }
}

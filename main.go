package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/jszroberto/kindle-words/kindledb"
	_ "github.com/mattn/go-sqlite3"
	// "log"
	"os"
)

func main() {

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Name = "kindle-words"
	app.Usage = "Provides methods to work with vocabulary builder "
	app.Version = "0.1.0"
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
			Name:   "delete-book",
			Usage:  "Delete book",
			Action: deleteBook,
		},
		{
			Name:   "cleanup",
			Usage:  "Cleanup",
			Action: cleanUp,
		},
		{
			Name:   "sync-with",
			Usage:  "Sync database with information from CSV",
			Action: syncDatabase,
		},
		{
			Name:   "delete-word",
			Usage:  "Delete word",
			Action: deleteWord,
		},
		{
			Name:   "books",
			Usage:  "show book title",
			Action: listBooks,
		},
	}

	app.Run(os.Args)

}

func listBooks(c *cli.Context) error {
	var language string
	if c.NArg() > 0 {
		language = c.Args()[0]
	}

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()

	books, err := kindle.ReadBooksInfo()
	if err != nil {
		return err
	}

	for _, book := range books {
		if language == "" || book.Language == language {
			fmt.Println(book.Title)
		}
	}
	return nil
}

func syncDatabase(c *cli.Context) error {
	var csvfile string
	if c.NArg() > 0 {
		csvfile = c.Args()[0]
	} else {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()
	words, err := readCSV(csvfile)

	if err != nil {
		return err
	}

	err = exportToCSV(csvfile+"-copy", words)
	return err
}

func cleanUp(c *cli.Context) error {

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()

	words, err := kindle.ReadWords()
	if err != nil {
		return err
	}

	frecuencies, err := getFrecuencies("en")
	if err != nil {
		return err
	}
	for _, word := range words {
		if _, ok := frecuencies[word.Value]; ok {
		} else {
			fmt.Println(word.Value)
			err := kindle.RemoveWord(word.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteBook(c *cli.Context) error {
	var bookName string
	if c.NArg() > 0 {
		bookName = c.Args()[0]
	} else {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()

	return kindle.RemoveBook(bookName)
}

func deleteWord(c *cli.Context) error {
	var word string
	if c.NArg() > 0 {
		word = c.Args()[0]
	} else {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()

	return kindle.RemoveWord(word)

}

func export(c *cli.Context) error {

	kindle, err := kindledb.NewKindleDB()
	if err != nil {
		return err
	}
	defer kindle.Close()

	words, err := kindle.ReadWords()
	if err != nil {
		return err
	}
	frecuencies, err := getFrecuencies("en")
	if err != nil {
		return err
	}

	words = sortWords(words, frecuencies)

	if c.Bool("html") {
		err = exportHtml("./export/html", words)
		if err != nil {
			return err
		}
	} else if c.Bool("vocab.com") {
		err = exportVocabularyCom("./export/vocab.com", words)
		if err != nil {
			return err
		}
	} else if c.Bool("evernote") {
		err = exportEvernote("./export/evernote", words)
		if err != nil {
			return err
		}
	} else if c.Bool("csv") {
		err = exportToCSV("result.csv", words)
		if err != nil {
			return err
		}
	} else {
		cli.ShowSubcommandHelp(c)
	}
	return nil
}

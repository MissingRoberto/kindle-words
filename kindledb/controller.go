package kindledb

import (
	"fmt"
	"github.com/kisielk/sqlstruct"
	"database/sql"
)

const (
	get_books_info = "SELECT * FROM BOOK_INFO"
	get_words      = "SELECT " +
		"w.word" +
		",group_concat(l.usage || '<i><small>' || b.title || '</i></small>' ,' <br/><br/> ') as usage " +
		",w.lang as language " +
		",b.title as book " +
		",w.category as mastered" +
		",count(l.usage) as count_usage " +
		"FROM " +
		"WORDS w " +
		"LEFT JOIN LOOKUPS l " +
		"on l.word_key=w.id " +
		"LEFT JOIN BOOK_INFO b " +
		"on b.guid=l.book_key " +
		"GROUP BY " +
		"w.word " +
		"ORDER BY book DESC, book DESC;"
	cleanup_lookups               = "DELETE FROM LOOKUPS WHERE NOT EXISTS(SELECT NULL FROM BOOK_INFO f WHERE f.id = book_key)"
	cleanup_words                 = "DELETE FROM WORDS WHERE id NOT IN (SELECT f.word_key FROM LOOKUPS f)"
	cleanup_lookups_when_no_words = "DELETE FROM LOOKUPS WHERE NOT EXISTS(SELECT NULL FROM WORDS f WHERE f.id = word_key)"
)

type KindleDB struct {
	db *sql.DB
}

func NewKindleDB() (*KindleDB,error) {
	path:= "/Volumes/Kindle/system/vocabulary/vocab.db"
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return &KindleDB{},err
	}
	return &KindleDB{db},nil
}

func (k *KindleDB) Close(){
	k.db.Close()
}

func (k *KindleDB) RemoveBook(title string) error {
	query := "DELETE FROM BOOK_INFO WHERE LOWER(title) LIKE LOWER('%" + title + "%')"
	output, err := k.db.Exec(query)

	if err != nil {
		return err
	}
	fmt.Println(output)
	_, err = k.db.Exec(cleanup_lookups)

	if err != nil {
		return err
	}
	_, err = k.db.Exec(cleanup_words)

	return err

}

func (k *KindleDB) RemoveWord(title string) error {
	query := "DELETE FROM WORDS WHERE  id  LIKE '%:" + title + "'"
	output, err := k.db.Exec(query)

	if err != nil {
		return err
	}
	fmt.Println(output)
	output, err = k.db.Exec(cleanup_lookups_when_no_words)
	fmt.Println(output)

	return err

}

func (k *KindleDB) ReadBooksInfo() ([]BookInfo, error) {
	objs := []BookInfo{}
	rows, err := k.db.Query(get_books_info)
	if err != nil {
		return []BookInfo{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var obj BookInfo
		err = sqlstruct.Scan(&obj, rows)
		if err != nil {
			return []BookInfo{}, err
		}
		objs = append(objs, obj)
	}

	err = rows.Err()
	if err != nil {
		return []BookInfo{}, err
	}
	return objs, nil
}

func (k *KindleDB) ReadWords() ([]Word, error) {
	words := []Word{}
	rows, err := k.db.Query(get_words)
	if err != nil {
		return []Word{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var word Word
		err = sqlstruct.Scan(&word, rows)
		if err != nil {
			return []Word{}, err
		}
		words = append(words, word)
	}

	err = rows.Err()
	if err != nil {
		return []Word{}, err
	}
	return words, nil
}

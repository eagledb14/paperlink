package dictionary

import (
	"fmt"
	"html"

	"github.com/eagledb14/paperlink/engagement"
	_ "modernc.org/sqlite"
)

type Dictionary struct {
	db *engagement.DbWrapper
}

type Word struct {
	Key int
	Word string
	Definition string
}

func LoadDictionary() Dictionary {
	db, err := engagement.Open("./dictionary.db?_busy_timeout=10000")
	if err != nil {
		panic("Unable to Create Dictionary")
	}

	createDictionaryTable(db)

	return Dictionary{
		db: db,
	}
}

func createDictionaryTable(db *engagement.DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS dictionary(
key INTEGER PRIMARY KEY,
word TEXT,
definition TEXT
)`)
}

func (d *Dictionary) InsertWord(word string, definition string) error {
	return d.db.Exec(`INSERT INTO dictionary(
word,
definition
) VALUES (?, ?)`, 
		html.EscapeString(word),
		html.EscapeString(definition))
}

func (d *Dictionary) Update(key int, word string, definition string) error {
	return d.db.Exec(`UPDATE dictionary SET word = ?, definition = ? WHERE key = ?`, 
		html.EscapeString(html.UnescapeString(word)),
		html.EscapeString(html.UnescapeString(definition)),
		key)
}

func (d *Dictionary) GetWords() []Word {
	rows, err := d.db.Query(`SELECT key, word, definition FROM dictionary`)
	if err != nil {
		return []Word{}
	}
	defer rows.Close()

	words := []Word{}
	for rows.Next() {
		newWord := Word{}
		if err := rows.Scan(&newWord.Key, &newWord.Word, &newWord.Definition); err != nil {
			continue
		}
		words = append(words, newWord)
	}

	return words
}

func (d *Dictionary) GetWord(key int) Word {
	row := d.db.QueryRow(`SELECT key, word, definition FROM dictionary WHERE key = ?`, key)
	newWord := Word{}

	if err := row.Scan(&newWord.Key, &newWord.Word, &newWord.Definition); err != nil {
		if err.Error() != "sql: no rows in result set" {
			fmt.Println("GetAsset: ", err)
		}
	}

	return newWord
}

func (d *Dictionary) GetWordFromWord(word string) (Word, error) {
	row := d.db.QueryRow(`SELECT key, word, definition FROM dictionary WHERE word = ?`, word)
	newWord := Word{}

	if err := row.Scan(&newWord.Key, &newWord.Word, &newWord.Definition); err != nil {
		return Word{}, err
	}

	return newWord, nil
}

func (d *Dictionary) Delete(key int) error {
	return d.db.Exec(`DELETE FROM dictionary WHERE key = ?`, key)
}



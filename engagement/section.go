package engagement

import "strings"


type Section struct {
	Key int
	Title string
	Body string
}

func createSectionTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS sections(
key INTEGER PRIMARY KEY,
"index" INTEGER,
title TEXT,
body TEXT
)`)
}

func (e *Engagement) InsertSection(title string, body string) error {
	return e.db.Exec(`INSERT INTO sections(
title,
body
) VALUES (?, ?)`, title, body)
}

func (e *Engagement) UpdateSection(key int, title string, body string) error {
	strings.ReplaceAll(body, "`", "'")
	return e.db.Exec(`UPDATE sections SET  title = ?, body = ? WHERE "key" = ?`, title, body, key)
}

func (e* Engagement) GetSections() []Section {
	rows, err := e.db.Query(`SELECT key, title, body FROM sections`)
	if err != nil {
		return []Section{}
	}
	defer rows.Close()

	sections := []Section{}
	for rows.Next() {
		newSection := Section{}
		if err := rows.Scan(&newSection.Key, &newSection.Title, &newSection.Body); err != nil {
			continue
		}
		sections = append(sections, newSection)
	}

	return sections
}

func (e *Engagement) GetSection(key string) Section {
	row := e.db.QueryRow(`SELECT key, title, body FROM sections WHERE key = ?`, key)
	newSection := Section{}
	row.Scan(&newSection.Key, &newSection.Title, &newSection.Body)

	return newSection
}

func (e *Engagement) DeleteSection(key int) error {
	return e.db.Exec(`DELETE FROM sections WHERE key = ?`, key)
}

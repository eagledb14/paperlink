package engagement

import (
	"sort"
	"strings"
	"html"
	"github.com/microcosm-cc/bluemonday"
)


type Section struct {
	Key int
	Index int
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
	var rowCount int
	err := e.db.QueryRow(`SELECT "index" FROM sections DESC LIMIT 1`).Scan(&rowCount)
	if err != nil {
		rowCount = 1
	}

	return e.db.Exec(`INSERT INTO sections(
title,
"index",
body
) VALUES (?, ?, ?)`,html.EscapeString(title), rowCount + 1, body)
}

func (e *Engagement) UpdateSection(key int, index int, title string, body string) error {
	strings.ReplaceAll(body, "`", "'")
	policy := bluemonday.UGCPolicy()
	policy.AllowStyles()
	body = policy.Sanitize(body)
	return e.db.Exec(`UPDATE sections SET "index" = ?, title = ?, body = ? WHERE "key" = ?`,
		index, html.EscapeString(html.UnescapeString(title)), body, key)
}

func (e* Engagement) GetSections() []Section {
	rows, err := e.db.Query(`SELECT key, "index", title, body FROM sections`)
	if err != nil {
		return []Section{}
	}
	defer rows.Close()

	sections := []Section{}
	for rows.Next() {
		newSection := Section{}
		if err := rows.Scan(&newSection.Key, &newSection.Index, &newSection.Title, &newSection.Body); err != nil {
			continue
		}
		sections = append(sections, newSection)
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Index < sections[j].Index
	})

	return sections
}

func (e *Engagement) GetSection(key string) Section {
	row := e.db.QueryRow(`SELECT key, "index", title, body FROM sections WHERE key = ?`, key)
	newSection := Section{}
	row.Scan(&newSection.Key, &newSection.Index, &newSection.Title, &newSection.Body)

	return newSection
}

func (e *Engagement) GetSectionFromIndex(index int) Section {
	row := e.db.QueryRow(`SELECT key, "index", title, body FROM sections WHERE "index" = ?`, index)
	newSection := Section{}
	row.Scan(&newSection.Key, &newSection.Index, &newSection.Title, &newSection.Body)

	return newSection
}

func (e *Engagement) DeleteSection(key int) error {
	return e.db.Exec(`DELETE FROM sections WHERE key = ?`, key)
}

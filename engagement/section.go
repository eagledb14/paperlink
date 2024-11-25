package engagement


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

func (e *Engagement) InsertSection(index int, title string, body string) error {
	return e.db.Exec(`INSERT INTO sections(
"index",
title,
body
) VALUES (?, ?, ?)`, index, title, body)
}

func (e *Engagement) UpdateSection(key int, index int, title string, body string) error {
	return e.db.Exec(`UPDATE sections SET "index" = ?, title = ?, body = ? WHERE "key" = ?`, index, title, body, key)
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

	return sections
}

func (e *Engagement) DeleteSection(key int) error {
	return e.db.Exec(`DELETE FROM sections WHERE key = ?`, key)
}

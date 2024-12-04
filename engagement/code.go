package engagement

import "html"



type Code struct {
	Key int
	Code string
	Paste string
}

func createCodeTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABlE IF NOT EXISTS codes(
key INTEGER PRIMARY KEY,
code TEXT,
paste TEXT
)`)
}

func (e *Engagement) InsertCode(code string, paste string) error {
	return e.db.Exec(`INSERT INTO codes(
code,
paste
)VALUES (?, ?)`, html.EscapeString(code), html.EscapeString(paste))
}

func (e *Engagement) GetCodes() []Code {
	rows, err := e.db.Query(`SELECT key, code, paste FROM codes`)
	if err != nil {
		return []Code{}
	}
	defer rows.Close()

	codes := []Code{}
	for rows.Next() {
		newCode := Code{}
		if err := rows.Scan(&newCode.Key, &newCode.Code, &newCode.Paste); err != nil {
			continue
		}
		codes = append(codes, newCode)
	}

	return codes
}

func (e *Engagement) DeleteCode(key int) error {
	return e.db.Exec(`DELETE FROM codes WHERE key = ?`, key)
}

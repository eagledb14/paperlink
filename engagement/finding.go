package engagement

import (
	"fmt"
	"time"
	"html"
)

type Finding struct {
	Key int

	//info low medium high critical
	Severity int
	TimeStamp time.Time

	Title string
	Body string

	DictionaryKey int
	AssetKey int
}

func createFindingTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS findings(
key INTEGER PRIMARY KEY,
severity INTEGER,
timeStamp INTEGER,
title TEXT,
body TEXT,
dictionaryKey int,
assetKey int
)`)
}

func (e *Engagement) InsertFinding(severity int, timeStamp time.Time, title string, body string, dictionaryKey int, assetKey int) (int, error) {
	// policy := bluemonday.UGCPolicy()
	// policy.AllowStyles()
	// policy.AllowElements("img")
	// policy.AllowAttrs("src").OnElements("img")
	//
	// body = policy.Sanitize(body)

	return e.db.ExecIndex(`INSERT INTO findings(
severity,
timeStamp,
title,
body,
dictionaryKey,
assetKey
) VALUES (?, ?, ?, ?, ?, ?)`, severity, timeStamp.Unix(), html.EscapeString(title), body, dictionaryKey, assetKey)
}

func (e *Engagement) UpdateFinding(key int, severity int, timeStamp time.Time, title string, body string, dictionaryKey int, assetKey int) error {
	// policy := bluemonday.UGCPolicy()
	// policy.AllowStyles()
	// body = policy.Sanitize(body)
	return e.db.Exec(`UPDATE findings SET severity = ?, timeStamp = ?, title = ?,  body= ?, dictionaryKey = ?, assetKey = ? WHERE "key" = ?`, 
		severity, timeStamp.Unix(), 
		html.EscapeString(html.UnescapeString(title)), body, dictionaryKey, assetKey, key)
}


func (e* Engagement) GetFindings() []Finding {
	rows, err := e.db.Query(`SELECT key, severity, timeStamp, title, body, dictionaryKey, assetKey FROM findings`)
	if err != nil {
		fmt.Println("Get Findings", err)
		return []Finding{}
	}
	defer rows.Close()

	findings := []Finding{}
	for rows.Next() {
		newFinding := Finding{}
		var timeStampUnix int64
		if err := rows.Scan(&newFinding.Key, &newFinding.Severity, &timeStampUnix, &newFinding.Title, &newFinding.Body, &newFinding.DictionaryKey, &newFinding.AssetKey); err != nil {
			fmt.Println(fmt.Errorf("GetFindings: %w", err))
			continue
		}
		newFinding.TimeStamp = time.Unix(timeStampUnix, 0)

		findings = append(findings, newFinding)
	}

	return findings
}

func (e* Engagement) GetFinding(key int) Finding {
	row := e.db.QueryRow(`SELECT key, severity, timeStamp, title, body, dictionaryKey, assetKey FROM findings WHERE key = ?`, key)
	newFinding := Finding{}
	var timeStampUnix int64
	if err := row.Scan(&newFinding.Key, &newFinding.Severity, &timeStampUnix, &newFinding.Title, &newFinding.Body, &newFinding.DictionaryKey, &newFinding.AssetKey); err != nil {
		fmt.Println("GetFinding:", err)
		return newFinding
	}
	newFinding.TimeStamp = time.Unix(timeStampUnix, 0)
	return newFinding
}

func (e* Engagement) GetFindingsWithAsset(key int) []Finding {
	rows, err := e.db.Query(`SELECT key, severity, timeStamp, title, body, dictionaryKey, assetKey FROM findings WHERE assetKey = ?`, key)
	if err != nil {
		return []Finding{}
	}
	defer rows.Close()

	findings := []Finding{}
	for rows.Next() {
		newFinding := Finding{}
		var timeStampUnix int64
		if err := rows.Scan(&newFinding.Key, &newFinding.Severity, &timeStampUnix, &newFinding.Title, &newFinding.Body, &newFinding.DictionaryKey, &newFinding.AssetKey); err != nil {
			fmt.Println(fmt.Errorf("GetFindings: %w", err))
			continue
		}
		newFinding.TimeStamp = time.Unix(timeStampUnix, 0)

		findings = append(findings, newFinding)
	}

	return findings
}


func (e *Engagement) DeleteFinding(key int) error {
	return e.db.Exec(`DELETE FROM findings WHERE key = ?`, key)
}

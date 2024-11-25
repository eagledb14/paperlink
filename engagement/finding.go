package engagement

import (
	"fmt"
	"time"
)

type Finding struct {
	Key int
	Severity int
	Title string
	StartDate time.Time
	Summary string
	Description string
}

func createFindingTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS findings(
key INTEGER PRIMARY KEY,
severity INTEGER,
title TEXT,
startDate INTEGER,
summary TEXT,
description TEXT
)`)
}

func (e *Engagement) InsertFinding(severity int, title string, startDate time.Time, summary string, description string) error {
	return e.db.Exec(`INSERT INTO findings(
severity,
title,
startDate,
summary,
description
) VALUES (?, ?, ?, ?, ?)`, severity, title, startDate.Unix(), summary, description)
}

func (e *Engagement) UpdateFinding(key int, severity int, title string, startDate time.Time, summary string, description string) error {
	return e.db.Exec(`UPDATE findings SET severity = ?, title = ?, startDate = ?, summary = ?, description = ? WHERE "key" = ?`, severity, title, startDate.Unix(), summary, description, key)
}

func (e* Engagement) GetFindings() []Finding {
	rows, err := e.db.Query(`SELECT key, severity, title, startDate, summary, description FROM findings`)
	if err != nil {
		fmt.Println(err)
		return []Finding{}
	}
	defer rows.Close()

	findings := []Finding{}
	for rows.Next() {
		newFinding := Finding{}
		var startDateUnix int64
		if err := rows.Scan(&newFinding.Key, &newFinding.Severity, &newFinding.Title, &startDateUnix, &newFinding.Summary, &newFinding.Description); err != nil {
			fmt.Println(fmt.Errorf("what %w", err))
			continue
		}
		newFinding.StartDate = time.Unix(startDateUnix, 0)

		fmt.Println(err)
		findings = append(findings, newFinding)
	}

	return findings
}

func (e *Engagement) DeleteFinding(key int) error {
	return e.db.Exec(`DELETE FROM findings WHERE key = ?`, key)
}

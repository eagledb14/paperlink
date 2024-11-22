package engagement

import (
	"time"
)

type Finding struct {
	key int
	severity int
	title string
	startDate time.Time
	summary string
	description string
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

func insertFinding(db *DbWrapper, severity int, title string, startDate time.Time, summary string, description string) error {
	return db.Exec(`INSERT INTO findings(
severity,
title,
startDate,
summary,
description
) VALUES (?, ?, ?, ?, ?)`, severity, title, startDate.Unix(), summary, description)
}


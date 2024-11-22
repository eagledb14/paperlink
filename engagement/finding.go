package engagement

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type FindingDb struct {
	db *sql.DB
}

func NewFindingDb(engagementName string, folderPath string) (error, FindingDb) {
	filePath := folderPath + engagementName + "/findings.db"

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return fmt.Errorf("FindingDb creation: %w", err), FindingDb{}
	}

	newFindingDb := FindingDb{
		db: db,
	}

	newFindingDb.createTable()
	return nil, newFindingDb
}

func (s *FindingDb) createTable() {
	tx, _ := s.db.Begin()
	tx.Exec(`CREATE TABLE IF NOT EXISTS sections(
key INTEGER PRIMARY KEY,
severity INTEGER,
title TEXT,
startDate TEXT,
)`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}

type Finding struct {
	key int
	severity int
	title string
	startDate time.Time
}

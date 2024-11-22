package engagement

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SectionDb struct {
	db *sql.DB
}

func NewSectionDb(engagementName string, folderPath string) (error, SectionDb) {
	filePath := folderPath + engagementName + "/sections.db"

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return fmt.Errorf("SectionDb creation: %w", err), SectionDb{}
	}

	newSectionDb := SectionDb{
		db: db,
	}

	newSectionDb.createTable()
	return nil, newSectionDb
}

func (s *SectionDb) createTable() {
	tx, _ := s.db.Begin()
	tx.Exec(`CREATE TABLE IF NOT EXISTS sections(
order INTEGER PRIMARY KEY,
title TEXT,
body TEXT,
)`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}

type Section struct {
	order int
	title string
	body string
}

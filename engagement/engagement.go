package engagement

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type EngagementDb struct {
	db *sql.DB
	folderPath string
}

func NewEngagementDb() EngagementDb {
	folderPath := "./engagements/"

	// Create the folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0700)
		if err != nil {
			panic("Could not find folder " + folderPath)
		}
	}

	db, err := sql.Open("sqlite", folderPath+"engagements.db")
	if err != nil {
		panic("Missing Resoruces")
	}

	newEngagementDb := EngagementDb{
		db: db,
		folderPath: folderPath,
	}

	newEngagementDb.createTable()
	return newEngagementDb
}

func NewTemplateDb() EngagementDb {
	folderPath := "./templates/"

	// Create the folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0700)
		if err != nil {
			panic("Could not find folder " + folderPath)
		}
	}

	db, err := sql.Open("sqlite", folderPath+"templates.db")
	if err != nil {
		panic("Missing Resoruces")
	}

	newEngagementDb := EngagementDb{
		db: db,
		folderPath: folderPath,
	}

	newEngagementDb.createTable()
	return newEngagementDb
}

func (e *EngagementDb) createTable() {
	tx, _ := e.db.Begin()
	tx.Exec(`CREATE TABLE IF NOT EXISTS engagements(
name TEXT PRIMARY KEY,
contact TEXT,
email TEXT,
)`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}

func (e *EngagementDb) insert(name string, contact string, email string) error {
	tx, _ := e.db.Begin()

	_, err := tx.Exec(`INSERT INTO engagements(
name,
contact,
email,
) VALUES (?,?,?)`, name, contact, email)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (e *EngagementDb) CreateEngagement(name string, contact string, email string) error {
	// create entry
	e.insert(name, contact, email)

	// create folder
	err := os.Mkdir(e.folderPath+name, 0700)

	// create section
	err, _ = NewSectionDb(name, e.folderPath)
	if err != nil {
		return fmt.Errorf("CreateBlankEngagement: %w", err)
	}

	// create asset
	err, _ = NewAssetDb(name, e.folderPath)
	if err != nil {
		return fmt.Errorf("CreateBlankEngagement: %w", err)
	}

	// create findings
	err, _ = NewFindingDb(name, e.folderPath)
	if err != nil {
		return fmt.Errorf("CreateBlankEngagement: %w", err)
	}

	return nil
}

func (e *EngagementDb) CreateEngagementFromTemplate(name string, templateName string, contact string, email string) error {
	copyPath := "./templates/"

	// create entry
	e.insert(templateName, contact, email)

	fsys := os.DirFS(copyPath + templateName)
	err := os.CopyFS(e.folderPath + name, fsys)
	if err != nil {
		return fmt.Errorf("CreateEngagementFromTemplate Copy: %w", err)
	}
	// copy section
	// copy asset
	// copy findings

	return nil
}

func (e *EngagementDb) DeleteEngagement(name string) error {
	err := os.RemoveAll(e.folderPath + name)

	tx, err := e.db.Begin()
	if err != nil {
		return fmt.Errorf("DeleteEngagementDb: %w", err)
	}

	_, nil := tx.Exec(`DELETE FROM engagements WHERE name = ?`, name)
	if err != nil {
		return fmt.Errorf("DeleteEngagementDb Deletion: %w", err)
	}

	return nil
}

type Engagement struct {
	name    string
	contact string
	email   string
}

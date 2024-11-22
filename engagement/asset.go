package engagement

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type AssetDb struct {
	db *sql.DB
}

func NewAssetDb(engagementName string, folderPath string) (error, AssetDb) {
	filePath := folderPath + engagementName + "/assets.db"

	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return fmt.Errorf("SectionDb creation: %w", err), AssetDb{}
	}

	newAssetDb := AssetDb{
		db: db,
	}

	newAssetDb.createTable()
	return nil, newAssetDb
}

func (a *AssetDb) createTable() {
	tx, _ := a.db.Begin()
	tx.Exec(`CREATE TABLE IF NOT EXISTS sections(
key INTEGER PRIMARY KEY,
parent TEXT,
name TEXT,
assetType TEXT,
)`)

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
	}
}

type Asset struct {
	key int
	parent string
	name string
	assetType string
}

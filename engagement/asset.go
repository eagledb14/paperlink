package engagement

import (
)

type Asset struct {
	key int
	parent string
	name string
	assetType string
}

func createAssetTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS assets(
key INTEGER PRIMARY KEY,
parent TEXT,
name TEXT,
assetType TEXT
)`)
}

func insertAsset(db *DbWrapper, parent string, name string, assetType string) error {
	return db.Exec(`INSERT INTO assets(
parent,
name,
assetType
) VALUES (?, ?, ?)`, parent, name, assetType)
}


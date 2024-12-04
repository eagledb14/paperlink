package engagement

import "fmt"

type Asset struct {
	Key int
	Parent string
	Name string
	AssetType string
}

func createAssetTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS assets(
key INTEGER PRIMARY KEY,
parent TEXT,
name TEXT,
assetType TEXT
)`)
}

func (e *Engagement) InsertAsset(parent string, name string, assetType string) error {
	return e.db.Exec(`INSERT INTO assets(
parent,
name,
assetType
) VALUES (?, ?, ?)`, parent, name, assetType)
}

func (e *Engagement) UpdateAsset(key int, parent string, name string, assetType string) error {
	return e.db.Exec(`UPDATE assets SET parent = ?, name = ?, assetType = ? WHERE "key" = ?`, parent, name, assetType, key)
}

func (e* Engagement) GetAssets() []Asset {
	rows, err := e.db.Query(`SELECT key, Parent, Name, AssetType FROM assets`)
	if err != nil {
		return []Asset{}
	}
	defer rows.Close()

	assets := []Asset{}
	for rows.Next() {
		newAsset := Asset{}
		if err := rows.Scan(&newAsset.Key, &newAsset.Parent, &newAsset.Name, &newAsset.AssetType); err != nil {
			continue
		}
		assets = append(assets, newAsset)
	}

	return assets
}

func (e* Engagement) GetAsset(key int) Asset {
	row := e.db.QueryRow(`SELECT key, parent, name, assetType FROM assets WHERE key = ?`, key)
	newAsset := Asset{}
	if err := row.Scan(&newAsset.Key, &newAsset.Parent, &newAsset.Name, &newAsset.AssetType); err != nil {
		fmt.Println("GetAsset: ", err)
	}

	return newAsset
}

func (e *Engagement) DeleteAsset(key int) error {
	return e.db.Exec(`DELETE FROM assets WHERE key = ?`, key)
}

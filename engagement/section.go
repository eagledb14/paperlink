package engagement

type Section struct {
	order int
	title string
	body string
}

func createSectionTable(db *DbWrapper) error {
	return db.Exec(`CREATE TABLE IF NOT EXISTS sections(
"order" INTEGER PRIMARY KEY,
title TEXT,
body TEXT
)`)
}

func insertSection(db *DbWrapper, order int, title string, body string) error {
	return db.Exec(`INSERT INTO sections(
"order",
title,
body
) VALUES (?, ?, ?)`, order, title, body)
}


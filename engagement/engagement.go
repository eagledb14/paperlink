package engagement

import (
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Engagement struct {
	db *DbWrapper
	folderPath string

	name    string
	contact string
	email   string
}

func NewEngagement(name string, contact string, email string) Engagement {
	folderPath := "./engagements/"
	return newDb(folderPath, name, contact, email)
}

func NewEngagementFromTemplate(templateName string, name string, contact string, email string) Engagement {
	copyPath := "./templates/" + templateName + ".db"
	destPath := "./engagements/" + name + ".db"

	// Create the folder if it doesn't exist
	if _, err := os.Stat("./engagements"); os.IsNotExist(err) {
		err := os.Mkdir("./engagements", 0700)
		if err != nil {
			panic("Could not find folder " + "./engagements")
		}
	}
	//todo: need to add dropping the engagement table, so the name doesn't show up from the template

	// copy template database
	err := copy(copyPath, destPath)
	if err != nil {
		fmt.Println(fmt.Errorf("CreateEngagementFromTemplate Copy: %w", err))
	}

	newEngagement := NewEngagement(name, contact, email)

	return newEngagement
}


func NewTemplate(name string) Engagement {
	folderPath := "./templates/"
	return newDb(folderPath, name, "", "")
}

func newDb(folderPath string, name string, contact string, email string) Engagement {
	// Create the folder if it doesn't exist
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0700)
		if err != nil {
			panic("Could not find folder " + folderPath)
		}
	}

	db, err := Open(folderPath+name+".db?_busy_timeout=10000")
	if err != nil {
		panic("Missing Resources")
	}

	newEngagement := Engagement{
		db: db,
		folderPath: folderPath,
		name: name,
		contact: contact,
		email: email,
	}

	newEngagement.createTable()

	// create engagement
	newEngagement.insertEngagement(name, contact, email)

	// create section
	err = createSectionTable(db)
	fmt.Println(err)

	// create asset
	createAssetTable(db)

	// create findings
	createFindingTable(db)

	return newEngagement
}

func (e *Engagement) createTable() {
	e.db.Exec(`CREATE TABLE IF NOT EXISTS engagements(
name TEXT PRIMARY KEY,
contact TEXT,
email TEXT
)`)
}

func (e *Engagement) insertEngagement(name string, contact string, email string) error {
	err := e.db.Exec(`INSERT INTO engagements(
name,
contact,
email
) VALUES (?,?,?)`, name, contact, email)
	return err
}

func (e *Engagement) Delete() {
	fmt.Println(e.folderPath + e.name + ".db")
	os.Remove(e.folderPath + e.name + ".db")
	e.Close()
}

func (e *Engagement) Close() {
	e.db.db.Close()
}


package engagement

import (
	"fmt"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

type Engagement struct {
	db *DbWrapper
	folderPath string

	Name    string
	Contact string
	Email   string
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

	// copy template database
	err := copy(copyPath, destPath)
	if err != nil {
		fmt.Println(fmt.Errorf("CreateEngagementFromTemplate Copy: %w", err))
	}

	newEngagement := NewEngagement(name, contact, email)
	newEngagement.deleteEngagement(templateName)
	if templateName == name {
		newEngagement.insertEngagement(name, contact, email)
	} 

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

	name = strings.TrimSpace(name)

	db, err := Open(folderPath+name+".db?_busy_timeout=10000")
	if err != nil {
		panic("Missing Resources")
	}

	newEngagement := Engagement{
		db: db,
		folderPath: folderPath,
		Name: name,
		Contact: contact,
		Email: email,
	}

	newEngagement.createTable()

	// create engagement
	newEngagement.insertEngagement(name, contact, email)

	// create section
	createSectionTable(db)

	// create asset
	createAssetTable(db)

	// create findings
	createFindingTable(db)

	return newEngagement
}

func LoadEngagements() []Engagement {
	files, err := os.ReadDir("./engagements/")
	if err != nil {
		return []Engagement{}
	}

	engagements := []Engagement{}
	for _, file := range files {
		if file.Type().IsRegular() {
			newEngagement, err := loadEngagement(file.Name(), "./engagements/")
			if err != nil {
				continue
			}
			engagements = append(engagements, newEngagement)
		}
	}
	return engagements
}

func LoadTemplates() []Engagement {
	files, err := os.ReadDir("./templates/")
	if err != nil {
		return []Engagement{}
	}

	engagements := []Engagement{}
	for _, file := range files {
		if file.Type().IsRegular() {
			newEngagement, err := loadEngagement(file.Name(), "./templates/")
			if err != nil {
				continue
			}
			engagements = append(engagements, newEngagement)
		}
	}

	return engagements
}

func loadEngagement(name string, folderPath string) (Engagement, error) {
	newEngagement := Engagement{}

	db, err := Open(folderPath+name+"?_busy_timeout=10000")
	if err != nil {
		panic("Missing Resources")
	}
	
	rows, err := db.Query(`SELECT name, contact, email FROM engagements`)
	if err != nil {
		fmt.Println(err)
		return Engagement{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&newEngagement.Name, &newEngagement.Contact, &newEngagement.Email)
		if err != nil {
			fmt.Println(err)
			return Engagement{}, err
		}

		newEngagement.db = db
		newEngagement.folderPath = folderPath

	}

	return newEngagement, nil
}

func (e *Engagement) createTable() error {
	return e.db.Exec(`CREATE TABLE IF NOT EXISTS engagements(
name TEXT PRIMARY KEY,
contact TEXT,
email TEXT
)`)
}

// use when creating engagement from template, so the template name doesn't stay
func (e *Engagement) deleteEngagement(templateName string) error {
	return e.db.Exec(`DELETE FROM engagements WHERE name = ?`, templateName)
}

func (e *Engagement) insertEngagement(name string, contact string, email string) error {
	err := e.db.Exec(`INSERT INTO engagements(
name,
contact,
email
) VALUES (?,?,?)`, name, contact, email)
	return err
}

func (e *Engagement) UpdateEngagement(name string, contact string, email string) error {
	copyPath := e.folderPath + e.Name + ".db"
	destPath :=  e.folderPath + name + ".db"

	e.db.mutex.Lock()

	e.Close()
	err := os.Rename(copyPath, destPath)
	if err != nil {
		e.db.mutex.Unlock()
		e.db, _ = Open(copyPath)
		return err
	}

	newDb, err := Open(destPath)
	if err != nil {
		e.db.mutex.Unlock()
		return err
	}
	e.db.mutex.Unlock()

	e.db.db = newDb.db
	err = e.db.Exec(`UPDATE engagements SET name = ?, contact = ?, email = ? WHERE name = ?`, name, contact, email, e.Name)

	e.Name = name
	e.Contact = contact
	e.Email = email

	return err
}

func (e *Engagement) Delete() {
	e.Close()
	os.Remove(e.folderPath + e.Name + ".db")
}

func (e *Engagement) Close() {
	e.db.Close()
}


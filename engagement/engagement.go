package engagement

import (
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"
	db "github.com/eagledb14/paperlink/db"
)

type Engagement struct {
	db *db.DbWrapper
	folderPath string

	Name    string
	Contact string
	Email   string
	TimeStamp time.Time
}

func NewEngagement(name string, contact string, email string) Engagement {
	folderPath := "./engagements/"
	return newDb(folderPath, name, contact, email)
}

func NewEngagementFromTemplate(templateName string, name string, contact string, email string) Engagement {
	name = html.EscapeString(name)
	templateName = html.EscapeString(templateName)
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
	err := db.Copy(copyPath, destPath)
	if err != nil {
		fmt.Println(fmt.Errorf("CreateEngagementFromTemplate Copy: %w", err))
	}

	newEngagement, err := loadEngagement(name + ".db", "./engagements/")
	if err != nil {
		panic(fmt.Errorf("Broken: %w: %s", err, newEngagement.Name))
	}
	newEngagement.deleteEngagement(templateName)

	newEngagement.Name = html.EscapeString(name)
	newEngagement.Contact = html.EscapeString(contact)
	newEngagement.Email = html.EscapeString(email)
	newEngagement.TimeStamp = time.Now()
	newEngagement.insertEngagement(name, contact, email, newEngagement.TimeStamp)

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

	db, err := db.Open(folderPath+name+".db?_busy_timeout=10000")
	if err != nil {
		panic("Unable to Create Engagement")
	}

	newEngagement := Engagement{
		db: db,
		folderPath: folderPath,
		Name: name,
		Contact: contact,
		Email: email,
		TimeStamp: time.Now(),
	}

	newEngagement.createTable()

	// create engagement
	err = newEngagement.insertEngagement(name, contact, email, newEngagement.TimeStamp)
	if err != nil {
		panic(fmt.Errorf("error making engagement table: %w", err))
	}

	// create section
	err = createSectionTable(db)
	if err != nil {
		panic(fmt.Errorf("error making section table: %w", err))
	}

	// create asset
	err = createAssetTable(db)
	if err != nil {
		panic(fmt.Errorf("error making asset table: %w", err))
	}

	// create findings
	err = createFindingTable(db)
	if err != nil {
		panic(fmt.Errorf("error making finding table: %w", err))
	}

	// create short codes
	err = createCodeTable(db)
	newEngagement.InsertCode("client name", name)
	newEngagement.InsertCode("client contact", contact)
	newEngagement.InsertCode("client email", email)
	if err != nil {
		panic(fmt.Errorf("error making code table: %w", err))
	}

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

	db, err := db.Open(folderPath+name+"?_busy_timeout=10000")
	if err != nil {
		panic("Missing Resources")
	}
	
	rows, err := db.Query(`SELECT name, contact, email, timeStamp FROM engagements`)
	if err != nil {
		fmt.Println(fmt.Errorf("load Engagement: %w", err))
		return Engagement{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var unix int64
		err = rows.Scan(&newEngagement.Name, &newEngagement.Contact, &newEngagement.Email, &unix)
		if err != nil {
			fmt.Println(err)
			return Engagement{}, err
		}

		newEngagement.TimeStamp = time.Unix(unix, 0)

		newEngagement.db = db
		newEngagement.folderPath = folderPath

	}

	return newEngagement, nil
}

func (e *Engagement) createTable() error {
	return e.db.Exec(`CREATE TABLE IF NOT EXISTS engagements(
name TEXT PRIMARY KEY,
contact TEXT,
email TEXT,
timeStamp INTEGER
)`)
}

// use when creating engagement from template, so the template name doesn't stay
func (e *Engagement) deleteEngagement(templateName string) error {
	return e.db.Exec(`DELETE FROM engagements WHERE name = ?`, templateName)
}

func (e *Engagement) insertEngagement(name string, contact string, email string, time time.Time) error {
	err := e.db.Exec(`INSERT INTO engagements(
name,
contact,
email,
timeStamp
) VALUES (?,?,?,?)`, name, contact, email, time.Unix())
	return err
}

func (e *Engagement) UpdateEngagement(name string, contact string, email string) error {
	copyPath := e.folderPath + e.Name + ".db"
	destPath :=  e.folderPath + name + ".db"

	e.db.Mutex.Lock()

	e.Close()
	err := os.Rename(copyPath, destPath)
	if err != nil {
		e.db.Mutex.Unlock()
		e.db, _ = db.Open(copyPath)
		return err
	}

	newDb, err := db.Open(destPath)
	if err != nil {
		e.db.Mutex.Unlock()
		return err
	}
	e.db.Mutex.Unlock()

	e.db.Db = newDb.Db
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


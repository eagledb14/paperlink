package engagement

import (
	"database/sql"
	"fmt"
	"os"
	"io"

	"sync"

	_ "modernc.org/sqlite"
)


type DbWrapper struct {
	db *sql.DB
	mutex sync.RWMutex
}

func Open(path string) (*DbWrapper, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return &DbWrapper{
		db: db,
	}, nil
}

func (d *DbWrapper) Exec(query string, args ...any) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	_, err := d.db.Exec(query, args...)
	return err
}

func (d *DbWrapper) ExecIndex(query string, args ...any) (int, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	result, err := d.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	// Get the ID of the last inserted row
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertID), nil
}

func (d *DbWrapper) Query(query string, args ...any) (*sql.Rows, error) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.db.Query(query, args...)
}

func (d *DbWrapper) QueryRow(query string, args ...any) *sql.Row {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return d.db.QueryRow(query, args...)
}

func (d *DbWrapper) Close() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.db.Close()
}

func copy(src, dst string) error {
	BUFFERSIZE := 200
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

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
	Db *sql.DB
	Mutex sync.RWMutex
}

func Open(path string) (*DbWrapper, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return &DbWrapper{
		Db: db,
	}, nil
}

func (d *DbWrapper) Exec(query string, args ...any) error {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	_, err := d.Db.Exec(query, args...)
	return err
}

func (d *DbWrapper) ExecIndex(query string, args ...any) (int, error) {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()

	result, err := d.Db.Exec(query, args...)
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
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.Db.Query(query, args...)
}

func (d *DbWrapper) QueryRow(query string, args ...any) *sql.Row {
	d.Mutex.RLock()
	defer d.Mutex.RUnlock()
	return d.Db.QueryRow(query, args...)
}

func (d *DbWrapper) Close() {
	d.Mutex.Lock()
	defer d.Mutex.Unlock()
	d.Db.Close()
}

func Copy(src, dst string) error {
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

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

func (d *DbWrapper) Query(query string) (*sql.Rows, error) {
	return d.db.Query(query)
}

func (d *DbWrapper) QueryRow(query string, args ...any) *sql.Row {
	return d.db.QueryRow(query, args...)
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

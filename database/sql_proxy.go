package database

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLProxy struct {
	sync.RWMutex

	databaseType string

	database *sql.DB

	state int
}

func (this *SQLProxy) Constructor() {
}
func (this *SQLProxy) Connect(databaseType, databasePath string) (err error) {
	this.databaseType = databaseType
	switch databaseType {
	case DATABASE_MYSQL:
		this.database, err = sql.Open("mysql", databasePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	case DATABASE_SQLITE3:
		this.database, err = sql.Open("sqlite3", databasePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = this.database.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
	this.database.SetConnMaxLifetime(time.Minute * 10)
	this.database.SetConnMaxIdleTime(time.Minute * 10)
	this.database.SetMaxOpenConns(1000)
	this.database.SetMaxIdleConns(1000)

	fmt.Println("创建数据库连接：", databaseType+"->"+databasePath, "::", this.database)
	this.state = STATE_CONNECTED

	go this.listen()

	return
}
func (this *SQLProxy) Prepare(query string) (*sql.Stmt, error) {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Prepare(query)
	case STATE_DISCONNECT:
	}
	return nil, errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) Query(query string, args ...any) (*sql.Rows, error) {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Query(query, args...)
	case STATE_DISCONNECT:
	}
	return nil, errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) QueryRow(query string, args ...any) *sql.Row {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.QueryRow(query, args...)
	case STATE_DISCONNECT:
	}
	return nil
}
func (this *SQLProxy) Exec(query string, args ...any) (sql.Result, error) {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Exec(query, args...)
	case STATE_DISCONNECT:
	}
	return nil, errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) Begin() (*sql.Tx, error) {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Begin()
	case STATE_DISCONNECT:
	}
	return nil, errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) Ping() error {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Ping()
	case STATE_DISCONNECT:
	}
	return errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) Close() error {
	switch this.state {
	case STATE_CONNECTED:
		return this.database.Close()
	case STATE_DISCONNECT:
	}
	return errors.New("数据库连接尚未建立...")
}
func (this *SQLProxy) listen() {
	for {
		if this.state == STATE_CONNECTED {
			var err error
			time.Sleep(time.Second * 10)
			err = this.database.Ping()
			if err != nil {
				this.state = STATE_DISCONNECT
			}
		}
	}
}

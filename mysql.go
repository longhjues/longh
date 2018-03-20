package longh

import (
	"database/sql"
	"log"
	"time"
)

// NewMySQL 创建一个连接MySQL的实体池
func NewMySQL(dbSource string, maxOpenConns, maxIdleConns int) (*sql.DB, error) {
	db, err := sql.Open("mysql", dbSource)
	if err != nil {
		return nil, err
	}
	if maxOpenConns > 0 {
		db.SetMaxOpenConns(maxOpenConns)
	}
	if maxIdleConns > 0 {
		db.SetMaxIdleConns(maxIdleConns)
	}

	go func() {
		for {
			err := db.Ping()
			if err != nil {
				log.Println("mysql db can't connect!")
			}
			time.Sleep(time.Minute)
		}
	}()
	log.Println("mysql connect successful")
	return db, nil
}

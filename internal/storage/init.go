package storage

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the SQLite database with the given schema
func InitDB(dbPath string, schemaPath string) *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	schema, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("스키마 파일 읽기 실패: %v", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		log.Fatalf("스키마 실행 실패: %v", err)
	}

	log.Println("SQLite DB 초기화 완료")
	return db
}

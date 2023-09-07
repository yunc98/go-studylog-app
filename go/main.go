package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "time/tzdata"

	_ "github.com/go-sql-driver/mysql"
)

var (
	user = os.Getenv("DB_USER")
	pass = os.Getenv("DB_PASS")
	dbname = os.Getenv("DB_NAME")
)

func main() {
	dsn := fmt.Sprintf("%s:%s@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, dbname)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	// SubjectをNewSubjectを使って作成
	s := NewSubject(db)

	// StudyLogをNewStudyLogを使って作成
	sl := NewStudyLog(db)

	// subjectsテーブルを作成
	if err := s.CreateSubjectsTable(); err != nil {
		log.Fatal(err)
	}

	// logsテーブルを作成
	if err := sl.CreateLogsTable(); err != nil {
		log.Fatal(err)
	}

	// HandlersをNewHandlersを使って作成
	hs := NewHandlers(sl)

	// ハンドラの登録
	http.HandleFunc("/", hs.ListHandler)
	http.HandleFunc("/save", hs.SaveHandler)
	http.HandleFunc("/summary", hs.SummaryHandler)
	
	fmt.Println("http://localhost:8080 で起動中...")
	// HTTPサーバーを起動する
	log.Fatal(http.ListenAndServe(":8080", nil))
}
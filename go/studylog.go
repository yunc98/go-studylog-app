package main

import (
	"database/sql"
)

type Log struct {
	ID       int
	Subject  string
	Duration int
}

// StudyLogの処理を行う型
type StudyLog struct {
	db *sql.DB
}

// 新しいStudyLogを作成する
func NewStudyLog(db *sql.DB) *StudyLog {
	// StudyLogのポインタを返す
	return &StudyLog{db: db}
}

// テーブルがなかったら作成する
func (sl *StudyLog) CreateTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS logs(
		id INTEGER PRIMARY KEY,
		subject TEXT NOT NULL,
		duration INTEGER NOT NULL
	);`

	_, err := sl.db.Exec((sqlStr))
	if err != nil {
		return err
	}

	return nil
}



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


// 新しいデータベースにLogを追加する
func (sl *StudyLog) AddLog(log *Log) error {
	const sqlStr = `INSERT INTO logs(subject, duration) VALUR (?,?);`

	_, err := sl.db.Exec(sqlStr, log.Subject, log.Duration)
	if err != nil {
		return err
	}

	return nil
}

// 最近追加したものを最大limit件だけ取得する
// エラーが発生したら第2戻り値で返す
func (sl *StudyLog) GetLogs(limit int) ([]*Log, error) {
	const sqlStr = `SELECT * FROM logs`

	rows, err := sl.db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var logs []*Log
	// 取得した各rowをLog型の変数にスキャンする
	for rows.Next() {
		var log Log
		err := rows.Scan(&log.ID, &log.Subject, &log.Duration)
		if err != nil {
			return nil, err
		}
		// logsスライスにスキャンしたrowsを追加する
		logs = append(logs, &log)
	}

	// rows.Next()のループ中にエラーが発生したかチェックする
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 取得したlogsスライスを返す
	return logs, nil
}

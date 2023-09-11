package main

import (
	"database/sql"
	"time"
)

type Log struct {
	ID       int
	SubjectId  int
	Duration int
	CreatedAt time.Time
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
func (sl *StudyLog) CreateLogsTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS logs(
		id INT AUTO_INCREMENT PRIMARY KEY,
		subject VARCHAR(255) NOT NULL,
		duration INT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := sl.db.Exec((sqlStr))
	if err != nil {
		return err
	}

	return nil
}


// 新しいデータベースにLogを追加する
func (sl *StudyLog) AddLog(log *Log) error {
	const sqlStr = `INSERT INTO logs(subject, duration) VALUES (?,?);`

	_, err := sl.db.Exec(sqlStr, log.SubjectId, log.Duration)
	if err != nil {
		return err
	}

	return nil
}

// 最近追加したものを最大limit件だけ取得する
// エラーが発生したら第2戻り値で返す
// できれば、件数ではなく月次の制限を設けたSQL文にしたい -> とりあえす元々の仕様で
func (sl *StudyLog) GetLogs(limit int) ([]*Log, error) {
	const sqlStr = `SELECT * FROM logs LIMIT ?`

	rows, err := sl.db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var logs []*Log
	// 取得した各rowをLog型の変数にスキャンする
	for rows.Next() {
		var log Log

		err := rows.Scan(&log.ID, &log.SubjectId, &log.Duration, &log.CreatedAt)
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

type Summary struct {
	Subject string
	Count int
	Sum int
}

// 集計結果を取得する
func (sl *StudyLog) GetSummaries() ([]*Summary, error) {
	const sqlStr = `
		SELECT
			subject,
			COUNT(1) as count,
			SUM(duration) as sum
		FROM
			logs
		GROUP BY
			subject
	`

	rows, err := sl.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var summaries []*Summary
	// 取得した各rowをSummary型の変数にスキャンする
	for rows.Next() {
		var s Summary

		err := rows.Scan(&s.Subject, &s.Count, &s.Sum)
		if err != nil {
			return nil, err
		}

		// summariesスライスにスキャンしたrowsを追加する
		summaries = append(summaries, &s)
	}

	// rows.Next()のループ中にエラーが発生したかチェックする
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 取得したsummariesスライスを返す
	return summaries, nil
}

// 平均を取得する
func (s *Summary) ComputeAvg() float64 {
	// Countが0だとゼロ除算になるため
	// そのまま0を返す
	if s.Count == 0 {
		return 0
	}

	return float64(s.Sum) / float64(s.Count)
}

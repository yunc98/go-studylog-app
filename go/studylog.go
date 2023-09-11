package main

import (
	"database/sql"
	"time"
)

type Log struct {
	ID       int
	SubjectId  int
	SubjectName string
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
		subjectId INT NOT NULL,
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
	const sqlStr = `INSERT INTO logs(subjectId, duration) VALUES (?,?);`

	_, err := sl.db.Exec(sqlStr, log.SubjectId, log.Duration)
	if err != nil {
		return err
	}

	return nil
}

// 最近追加したものを最大limit件だけ取得する
// エラーが発生したら第2戻り値で返す
func (sl *StudyLog) GetLogs(limit int) ([]*Log, error) {
	const sqlStr = `SELECT logs.id, logs.subjectId, subjects.subject, logs.duration
		FROM logs
		LEFT JOIN subjects ON logs.subjectId = subjects.id
		LIMIT ?`

	rows, err := sl.db.Query(sqlStr, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var logs []*Log
	// 取得した各rowをLog型の変数にスキャンする
	for rows.Next() {
		var log Log

		err := rows.Scan(&log.ID, &log.SubjectId, &log.SubjectName, &log.Duration)
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

type SummaryBySubject struct {
	SubjectName string
	Count int
	Sum int
}

// Subject毎の集計結果を取得する
func (sl *StudyLog) GetSummariesBySubject() ([]*SummaryBySubject, error) {
	const sqlStr = `
		SELECT
			subjects.subject,
			COUNT(1) as count,
			SUM(duration) as sum
		FROM
			logs
		LEFT JOIN subjects ON logs.subjectId = subjects.id
		GROUP BY logs.subjectId;`
	
	rows, err := sl.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される
	
	var summariesBySubject []*SummaryBySubject
	// 取得した各rowをSummaryBySubject型の変数にスキャンする
	for rows.Next() {
		var s SummaryBySubject
		
		err := rows.Scan(&s.SubjectName, &s.Count, &s.Sum)
		if err != nil {
			return nil, err
		}
		
		// summariesBySubjectスライスにスキャンしたrowsを追加する
		summariesBySubject = append(summariesBySubject, &s)
	}
	
	// rows.Next()のループ中にエラーが発生したかチェックする
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	// 取得したsummariesBySubjectスライスを返す
	return summariesBySubject, nil
}

// Subject毎の平均を取得する
func (s *SummaryBySubject) SubjectAvg() float64 {
	// Countが0だとゼロ除算になるため
	// そのまま0を返す
	if s.Count == 0 {
		return 0
	}
	
	return float64(s.Sum) / float64(s.Count)
}

type SummaryByMonth struct {
	Month string
	Count int
	Sum int
}

// Month毎の集計結果を取得する
func (sl *StudyLog) GetSummariesByMonth() ([]*SummaryByMonth, error) {
	const sqlStr = `
		SELECT
			DATE_FORMAT(logs.created_at, '%Y-%m') as month,
			COUNT(1) as count,
			SUM(duration) as sum
		FROM logs
		GROUP BY month;`

	rows, err := sl.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var summariesByMonth []*SummaryByMonth
	// 取得した各rowをSummaryByMonth型の変数にスキャンする
	for rows.Next() {
		var s SummaryByMonth

		err := rows.Scan(&s.Month, &s.Count, &s.Sum)
		if err != nil {
			return nil, err
		}

		// summariesByMonthスライスにスキャンしたrowsを追加する
		summariesByMonth = append(summariesByMonth, &s)
	}

	// rows.Next()のループ中にエラーが発生したかチェックする
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	// 取得したsummariesByMonthスライスを返す
	return summariesByMonth, nil
}

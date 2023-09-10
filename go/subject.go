package main

import (
	"database/sql"
	"time"
)

type SubjectItem struct {
	ID int
	Subject string
	CreatedAt time.Time
}

// Subjectの処理を行う型
type Subject struct {
	db *sql.DB
}

// 新しいSubjectを作成する
func NewSubject(db *sql.DB) *Subject {
	// Subjectのポインタを返す
	return &Subject{db: db}
}

// テーブルがなかったら作成する
func (s *Subject) CreateSubjectsTable() error {
	const sqlStr = `CREATE TABLE IF NOT EXISTS subjects(
		id INT AUTO_INCREMENT PRIMARY KEY,
		subject VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := s.db.Exec(sqlStr)
	if err != nil {
		return err
	}

	return nil
}

// 新しいデータベースにSubjectを追加する
func (s *Subject) AddSubject(si *SubjectItem) error {
	const sqlStr = `INSERT INTO subjects(subject) VALUE (?);`

	_, err := s.db.Exec(sqlStr, si.Subject)
	if err != nil {
		return err
	}

	return nil
}

// subjectsを取得する
func (s *Subject) GetSubjects() ([]*SubjectItem, error) {
	const sqlStr = `SELECT * FROM subjects`

	rows, err := s.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // 関数終了時にCloseが呼び出される

	var subjects []*SubjectItem
	// 取得した各rowをSubjectItem型の変数にスキャンする
	for rows.Next() {
		var subject SubjectItem

		err := rows.Scan(&subject.ID, &subject.Subject, &subject.CreatedAt)
		if err != nil {
			return nil, err
		}

		// subjectsスライスにスキャンしたrowsを追加する
		subjects = append(subjects, &subject)
	}

	// rows.Next()のループ中にエラーが発生したかチェックする
	if err = rows.Err(); err != nil {
		return nil, err
	}

	// 取得したsubjectsスライスを返す
	return subjects, nil
}
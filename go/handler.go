package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// HTTPハンドラを集めた型
type Handlers struct {
	sl *StudyLog
}

// Handlersを作成する
func NewHandlers(sl *StudyLog) *Handlers {
	return &Handlers{sl: sl}
}

// ListHandlerで使用するテンプレート
var listTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8"/>
			<title>Study Log📚<title>
		</head>
		<body>
			<h1>Study Log📚</h1>
			<h2>Add new log</h2>
			<form method="post" action="/save">
				<label for="subject">Subject</label>
				<input name="subject" type="text">
				<label for="duration">Duration</label>
				<input name="duration" type="number">
				<input type="submit" value="Add">
			</form>

			<h2>Latest logs : {{len .}}(<a href="/summary">Summary</a>)</h2>
			{{- if . -}}
			<table border="1">
				<tr><th>Subject</th><th>Duration</th></tr>
				{{- range .}}
				<tr><td>{{.Subject}}</td><td>{{.Duration}}</td></tr>
				{{- end}}
			</table>
			{{- else}}
				No record
			{{- end}}
		</body>
	</html>
`))

// 最新の入力データを表示するハンドラ
func (hs *Handlers) ListHandler(w http.ResponseWriter, r *http.Request) {
	// 最新の10件を取得する
	logs, err := hs.sl.GetLogs(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得したlogsをテンプレートに埋め込む
	if err := listTmpl.Execute(w, logs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// 保存するハンドラ
func (hs *Handlers) SaveHandler(w http.ResponseWriter, r *http.Request) {
	// リクエストがPOSTメソッドかチェックする
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed // 405
		http.Error(w, http.StatusText(code), code)
		return
	}

	// リクエストのフォームからsubjectフィールドの値を取得する
	subject := r.FormValue("subject")
	if subject == "" {
		// 空の文字列だったら、400を返す
		http.Error(w, "Subject not entered", http.StatusBadRequest)
		return
	}

	// リクエストのフォームからdurationフィールドの値を取得して、int(整数)に変換する
	duration, err := strconv.Atoi(r.FormValue("duration"))
	if err != nil {
		// intに変換できなかったら、400を返す
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 取得した値をもとに、新しいLogインスタンスを作成する
	log := &Log{
		Subject: subject,
		Duration: duration,
	}

	// logを保存する
	if err := hs.sl.AddLog(log); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	// logの保存に成功したら、ルートパスにリダイレクトする
	http.Redirect(w, r, "/", http.StatusFound) // 302
}

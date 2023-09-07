package main

import (
	"html/template"
	"net/http"
	"strconv"
)

// HTTPハンドラを集めた型
type Handlers struct {
	s  *Subject
	sl *StudyLog
}

// Handlersを作成する
func NewHandlers(s *Subject, sl *StudyLog) *Handlers {
	return &Handlers{s: s, sl: sl}
}

// ListHandlerで使用するテンプレート
var listTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8"/>
			<title>Study Log📚</title>
		</head>
		<body>
			<h1>Study Log📚</h1>
			<h2>Add new subject</h2>
			<form method="post" action="/save-subject">
				<label for="subject">Subject</label>
				<input name="subject" type="text">
				<input type="submit" value="Add">
			</form>

			<h2>Add new log</h2>
			<form method="post" action="/save-log">
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

// subjectを保存するハンドラ
func (hs *Handlers) SaveSubjectHandler(w http.ResponseWriter, r *http.Request) {
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

	// 取得した値をもとに、新しいSubjectItemインスタンスを作成する
	subjectItem := &SubjectItem{
		Subject: subject,
	}

	// subjectItemを保存する
	if err := hs.s.AddSubject(subjectItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// logの保存に成功したら、ルートパスにリダイレクトする
	http.Redirect(w, r, "/", http.StatusFound) // 302
}

// logを保存するハンドラ
func (hs *Handlers) SaveLogHandler(w http.ResponseWriter, r *http.Request) {
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

var summaryTmpl = template.Must(template.New("summary").Parse(`<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8"/>
			<title>Study Log📚 Summary</title>
			<script src="https://www.gstatic.com/charts/loader.js"></script>
			<script>
				google.charts.load('current', {'packages':['corechart']});
				google.charts.setOnLoadCallback(drawChart);

				function drawChart() {
				var data = google.visualization.arrayToDataTable([
					['Subject', 'Duration'],
					{{- range . -}}
					['{{js .Subject}}', {{.Sum}}],
					{{- end -}}
				]);
			
			var options = { title: 'Percentage' };
			var chart = new google.visualization.PieChart(document.getElementById('piechart'));
			chart.draw(data, options);
			}
			</script>
		</head>
		<body>
			<h1>Summary</h1>
			{{- if . -}}
			<div id="piechart" style="width:400px; height:300px;"></div>
			<table border="1">
				<tr><th>Subject</th><th>Total</th><th>Average</th></tr>
				{{- range .}}
				<tr><td>{{.Subject}}</td><td>{{.Sum}} hours</td><td>{{.ComputeAvg}} hours</tr>
				{{- end}}
			</table>
			{{- else}}
				No record
			{{- end}}

			<div><a href="/">Back</a></div>
		</body>
	</html>
`))

func (hs *Handlers) SummaryHandler(w http.ResponseWriter, r *http.Request) {
	// summariesを取得する
	summaries, err := hs.sl.GetSummaries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	// summariesをテンプレートに埋め込む
	if err := summaryTmpl.Execute(w, summaries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}
}

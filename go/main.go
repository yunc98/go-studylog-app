package main

import (
	"fmt"
)

type Log struct {
	Subject string
	Duration int
}

func main() {
	// 入力するデータの件数
	var n int
	fmt.Print("How many logs?>")
	fmt.Scan(&n)

	logs := make([]Log, 0, n)

	for i := 0; i < n; i++ {
		logs = inputLog(logs)
	}

	showLogs(logs)
}

func inputLog(logs []Log) []Log {
	var log Log

	fmt.Print("Subject>")
	fmt.Scan(&log.Subject)

	fmt.Print("Duration>")
	fmt.Scan(&log.Duration)

	logs = append(logs, log)

	return logs
}

func showLogs(logs []Log) {
	fmt.Println("----------")

	for i := 0; i < len(logs); i++ {
		fmt.Printf("Studied %s for %d hours\n", logs[i].Subject, logs[i].Duration)
	}

	fmt.Println("----------")
}
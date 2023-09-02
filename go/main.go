package main

import (
	"fmt"
)

type Log struct {
	Subject string
	Duration int
}

func main() {
	log := inputLog()

	fmt.Printf("Studied %s for %d hours\n", log.Subject, log.Duration)
}

func inputLog() Log {
	var log Log

	fmt.Print("Subject>")
	fmt.Scan(&log.Subject)

	fmt.Print("Duration>")
	fmt.Scan(&log.Duration)

	return log
}
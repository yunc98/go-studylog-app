package main

import "fmt"

func main() {
	var subject string
	var time int

	fmt.Print("Subject>")
	fmt.Scan(&subject)

	fmt.Print("Time>")
	fmt.Scan(&time)

	fmt.Printf("Studied %s for %d hours\n", subject, time)
}
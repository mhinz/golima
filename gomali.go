package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	foundIssue := false

	if file, err := os.Open("/data/github/vim-galore/README.md"); err != nil {
		log.Fatal(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			foundIssue = checkRules(scanner.Text()) || foundIssue
		}
		if err = scanner.Err(); err != nil {
			file.Close()
			log.Fatal(err)
		}
		file.Close()
	}

	if foundIssue {
		os.Exit(1)
	}
	os.Exit(0)
}

func checkRules(s string) bool {
	foundIssue := ruleTooLong(s)
	return foundIssue
}

func ruleTooLong(s string) bool {
	if len(s) > 80 {
		log.Printf("Too long: %s\n", s)
		return true
	} else {
		log.Printf("Len: %d\n", len(s))
	}
	return false
}

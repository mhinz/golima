package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	foundIssue := false
	log.SetFlags(log.Lshortfile)

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
		log.Println("Line longer than 80 characters.")
		return true
	}
	return false
}

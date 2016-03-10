package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

type Context struct {
	prevLine string
	curLine  string
	nextLine string
}

var (
	reHeader = regexp.MustCompile("^#{1,6}")
)

func main() {
	foundIssue := false
	log.SetFlags(log.Lshortfile)

	if file, err := os.Open("/data/github/vim-galore/README.md"); err != nil {
		log.Fatal(err)
	} else {
		scanner := bufio.NewScanner(file)
		ctx := Context{"", scanner.Text(), scanner.Text()}
		for scanner.Scan() {
			foundIssue = ctx.checkRules() || foundIssue
			ctx = Context{ctx.curLine, ctx.nextLine, scanner.Text()}
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

func (ctx *Context) checkRules() bool {
	foundIssue := false
	foundIssue = ctx.ruleLineLength() || foundIssue
	foundIssue = ctx.ruleProperHeader() || foundIssue
	return foundIssue
}

func (ctx *Context) ruleProperHeader() bool {
	if reHeader.MatchString(ctx.curLine) {
		if len(ctx.prevLine) > 0 || len(ctx.nextLine) > 0 {
			log.Println("Header must be surrounded by blank lines")
			return true
		}
	}
	return false
}

func (ctx *Context) ruleLineLength() bool {
	if len(ctx.curLine) > 80 {
		log.Println("Line longer than 80 characters.")
		return true
	}
	return false
}

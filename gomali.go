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
	foundIssue = 0
	reHeader   = regexp.MustCompile("^#{1,6}")
)

func main() {
	log.SetFlags(log.Lshortfile)

	if file, err := os.Open("/data/github/vim-galore/README.md"); err != nil {
		log.Fatal(err)
	} else {
		scanner := bufio.NewScanner(file)
		ctx := Context{"", scanner.Text(), scanner.Text()}
		for scanner.Scan() {
			ctx.checkRules()
			ctx = Context{ctx.curLine, ctx.nextLine, scanner.Text()}
		}
		if err = scanner.Err(); err != nil {
			file.Close()
			log.Fatal(err)
		}
		file.Close()
	}
	os.Exit(foundIssue)
}

func (ctx *Context) checkRules() {
	ctx.ruleLineLength()
	ctx.ruleProperHeader()
}

func (ctx *Context) ruleProperHeader() {
	if reHeader.MatchString(ctx.curLine) {
		if len(ctx.prevLine) > 0 || len(ctx.nextLine) > 0 {
			log.Println("Header must be surrounded by blank lines")
			foundIssue = 1
		}
	}
}

func (ctx *Context) ruleLineLength() {
	if len(ctx.curLine) > 80 {
		log.Println("Line longer than 80 characters.")
		foundIssue = 1
	}
}

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// The current line with one line of surrounding context.
type Context struct {
	filename  string
	curLineNr int
	curLine   string
	prevLine  string
	nextLine  string
}

var (
	foundIssue = 0  // process return value
	reHeader   = regexp.MustCompile("^#{1,6}")
)

func main() {
	log.SetFlags(log.Lshortfile)
	file, err := os.Open("/data/github/vim-galore/README.md")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ctx := Context{"foo.vim", -1, scanner.Text(), "", scanner.Text()}
	for scanner.Scan() {
		ctx.checkRules()
		ctx = Context{ctx.filename, ctx.curLineNr+1, ctx.nextLine, ctx.curLine, scanner.Text()}
	}
	if err = scanner.Err(); err != nil {
		file.Close()
		log.Fatal(err)
	}

	os.Exit(foundIssue)
}

func (ctx *Context) print(msg string) {
	fmt.Printf("%s:%d:%s\n", ctx.filename, ctx.curLineNr, msg)
}

// Check all available rules for the current context.
func (ctx *Context) checkRules() {
	ctx.ruleLineLength()
	ctx.ruleProperHeader()
}

// Check if a potential header is surrounded by blank lines.
func (ctx *Context) ruleProperHeader() {
	if reHeader.MatchString(ctx.curLine) {
		if len(ctx.prevLine) > 0 || len(ctx.nextLine) > 0 {
			ctx.print("Header must be surrounded by blank lines")
			foundIssue = 1
		}
	}
}

// Check if the current line is longer than 80 characters.
func (ctx *Context) ruleLineLength() {
	if len(ctx.curLine) > 80 {
		ctx.print("Line longer than 80 characters.")
		foundIssue = 1
	}
}

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
	reBullet   = regexp.MustCompile("^\\s*- ")
	reHeader   = regexp.MustCompile("^#{1,6} ")
	reLink     = regexp.MustCompile("[\\w+\\]\\([\\w#]+\\)")
)

func main() {
	log.SetFlags(log.Lshortfile)
	for _, filename := range os.Args[1:] {
		checkFile(filename)
	}
	os.Exit(foundIssue)
}

func checkFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	ctx := Context{filename, -1, scanner.Text(), "", scanner.Text()}
	for scanner.Scan() {
		ctx.checkRules()
		ctx = Context{ctx.filename, ctx.curLineNr+1, ctx.nextLine, ctx.curLine, scanner.Text()}
	}
	// Since we called scanner.Text() twice before scanner.Scan(),
	// we still have to check the last two lines at this point.
	for i := 0; i < 1; i++ {
		ctx.checkRules()
		ctx = Context{ctx.filename, ctx.curLineNr+1, ctx.nextLine, ctx.curLine, scanner.Text()}
	}
	if err = scanner.Err(); err != nil {
		file.Close()
		log.Fatal(err)
	}
}

func (ctx *Context) print(msg string) {
	fmt.Printf("%s:%d:%s\n", ctx.filename, ctx.curLineNr, msg)
}

// Check all available rules for the current context.
func (ctx *Context) checkRules() {
	if reBullet.MatchString(ctx.curLine) {
		return
	}
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
	if len(ctx.curLine) > 80 && !reLink.MatchString(ctx.curLine) {
		ctx.print("Line longer than 80 characters.")
		foundIssue = 1
	}
}

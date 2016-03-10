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
	filename    string
	curLineNr   int
	curLine     string
	prevLine    string
	nextLine    string
}

var (
	foundIssue  = 0  // process return value
	inCodeBlock = false
	inTable     = false
	reBullet    = regexp.MustCompile("^\\s*- ")
	reCodeBlock = regexp.MustCompile("^```")
	reHeader    = regexp.MustCompile("^#{1,6} ")
	reLink      = regexp.MustCompile("[\\w+\\]\\([\\w#]+\\)")
	reTable     = regexp.MustCompile("^\\|")
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
	foundIssue = 1
}

// Check all available rules for the current context.
func (ctx *Context) checkRules() {
	if reCodeBlock.MatchString(ctx.curLine) {
		ctx.ruleProperCodeBlock()
		return
	} else if inCodeBlock {
		return
	} else if reTable.MatchString(ctx.curLine) {
		ctx.ruleProperTable()
		return
	} else if reBullet.MatchString(ctx.curLine) {
		return
	}
	ctx.ruleConsecutiveBlankLines()
	ctx.ruleLineLength()
	ctx.ruleProperHeader()
}

// CHeck if there are consecutive blank lines.
func (ctx *Context) ruleConsecutiveBlankLines() {
	if ctx.prevLine == "" && ctx.curLine == "" && ctx.curLineNr > 1 {
		ctx.print("No reason for consecutive blank lines.")
	}
}

// Check if table is surrounded by blank lines.
func (ctx *Context) ruleProperTable() {
	if inTable {
		if ctx.nextLine == "" {
			inTable = false
		} else if !reTable.MatchString(ctx.nextLine) {
			ctx.print("Tables must be surrounded by blank lines.")
		}
	} else {
		if ctx.prevLine == "" {
			inTable = true
		} else {
			ctx.print("Tables must be surrounded by blank lines.")
		}
	}
}

// Check if code block is surrounded by blank lines.
func (ctx *Context) ruleProperCodeBlock() {
	if inCodeBlock {
		if len(ctx.nextLine) > 0 {
			ctx.print("Code blocks must be surrounded by blank lines.")
		}
		inCodeBlock = false
	} else {
		if len(ctx.prevLine) > 0 {
			ctx.print("Code blocks must be surrounded by blank lines.")
		}
		inCodeBlock = true
	}
}

// Check if a potential header is surrounded by blank lines.
func (ctx *Context) ruleProperHeader() {
	if reHeader.MatchString(ctx.curLine) {
		if len(ctx.prevLine) > 0 || len(ctx.nextLine) > 0 {
			ctx.print("Header must be surrounded by blank lines.")
		}
	}
}

// Check if the current line is longer than 80 characters.
func (ctx *Context) ruleLineLength() {
	if len(ctx.curLine) > 80 && !reLink.MatchString(ctx.curLine) {
		ctx.print("Line longer than 80 characters.")
	}
}

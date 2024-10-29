package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var uFlag bool
var dFlag bool
var cFlag bool
var iFlag bool
var sFlag int
var app = App{}

type App struct {
	Input              *os.File
	Output             *os.File
	PreviousLine       *string
	PreviousLineOrigin *string
}

func init() {
	flag.BoolVar(&cFlag, "c", false, "count")
	flag.BoolVar(&uFlag, "u", false, "count")
	flag.BoolVar(&dFlag, "d", false, "count")
	flag.BoolVar(&iFlag, "i", false, "count")
	flag.IntVar(&sFlag, "s", 0, "skip")
}

// Выполняет функционал первых 3х флагов
func cduArg(scanner *bufio.Scanner) {
	duplicates := 0
	switch true {
	case cFlag:
		fmt.Printf("cFlag: %v :%+v\n", cFlag, flag.Args())
		fmt.Printf("Посчитать строки\n")
		for scanner.Scan() {
			if iFlag {
				line := strings.ToLower(scanner.Text())
				lineOrigin := scanner.Text()
				if app.PreviousLine == nil {
					app.PreviousLine = &line
					app.PreviousLineOrigin = &lineOrigin
				} else if *app.PreviousLine == line {
					if duplicates == 1 {
						fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PreviousLineOrigin)
					}
					app.PreviousLine = &line
					app.PreviousLineOrigin = &lineOrigin
					duplicates++
				} else {
					app.PreviousLine = &line
					app.PreviousLineOrigin = &lineOrigin
					duplicates = 0
				}
				duplicates++
			} else {
				line := scanner.Text()
				if app.PreviousLine == nil {
					app.PreviousLine = &line
				} else if *app.PreviousLine != line {
					fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PreviousLine)
					app.PreviousLine = &line
					duplicates = 0
				}
				duplicates++
			}
		}
		if app.PreviousLine != nil {
			fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PreviousLine)
		}
	case uFlag:

		fmt.Printf("uFlag: %v :%+v\n", uFlag, flag.Args())
		fmt.Printf("Вывод не повторяющиеся строк\n")
		for scanner.Scan() {
			line := scanner.Text()
			if app.PreviousLine == nil {
				app.PreviousLine = &line
			} else if *app.PreviousLine != line {
				if duplicates > 1 {
					fmt.Fprintf(app.Output, "%s\n", line)
				}
				app.PreviousLine = &line
				duplicates = 0
			}
			duplicates++
		}
		if app.PreviousLine != nil {
			if duplicates == 1 {
				fmt.Fprintf(app.Output, "%s\n", *app.PreviousLine)
			}
		}
	case dFlag:
		fmt.Printf("dFlag: %v :%+v\n", dFlag, flag.Args())
		fmt.Printf("Вывод повторяющиеся строки\n")
		for scanner.Scan() {
			line := scanner.Text()
			if app.PreviousLine == nil {
				app.PreviousLine = &line
			} else if *app.PreviousLine != line {
				if duplicates > 1 {
					fmt.Fprintf(app.Output, "%s\n", *app.PreviousLine)
				}
				app.PreviousLine = &line
				duplicates = 0
			}
			duplicates++
		}
		if app.PreviousLine != nil {
			if duplicates != 1 {
				fmt.Fprintf(app.Output, "%s\n", *app.PreviousLine)
			}
		}
	}
}

func openFile(filename string) *os.File {
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	return fh
}

func main() {
	flag.Parse()

	argN := len(flag.Args())
	switch argN {
	case 0:
		app.Input = os.Stdin
		app.Output = os.Stdout
	case 1:
		app.Input = openFile(flag.Args()[0])
		defer app.Input.Close()
		app.Output = os.Stdout
	case 2:
		app.Input = openFile(flag.Args()[0])
		app.Output = openFile(flag.Args()[1])
		defer app.Input.Close()
		defer app.Output.Close()
	default:
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	stage := bufio.NewScanner(app.Input)
	cduArg(stage)
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	//"strings"
)

var uFlag bool
var dFlag bool
var cFlag bool
var iFlag bool
var sFlag int

func init() {
	flag.BoolVar(&cFlag, "c", false, "count")
	flag.BoolVar(&uFlag, "u", false, "count")
	flag.BoolVar(&dFlag, "d", false, "count")
	flag.BoolVar(&iFlag, "i", false, "count")
	flag.IntVar(&sFlag, "s", 0, "skip")
}

type App struct {
	Input        *os.File
	Output       *os.File
	PreviousLine *string
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

	app := App{}

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
	duplicates := 0

	switch true {
	case cFlag:
		fmt.Printf("\ncFlag: %v :%+v\n", cFlag, flag.Args())
		fmt.Printf("\nПосчитать строки\n\n")
		for stage.Scan() {
			line := stage.Text()
			if app.PreviousLine == nil {
				app.PreviousLine = &line
			} else if *app.PreviousLine != line {
				if cFlag {
					fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PreviousLine)
				}
				app.PreviousLine = &line
				duplicates = 0
			}
			duplicates++
		}
		if app.PreviousLine != nil {
			if cFlag {
				fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PreviousLine)
			}
		}
	case uFlag:
		fmt.Println("Вывод не повторяющиеся строк")

	case dFlag:
		fmt.Println("Вывести повторяющиеся строки")

	default:
		for stage.Scan() {
			line := stage.Text()
			if app.PreviousLine == nil {
				app.PreviousLine = &line
			} else if *app.PreviousLine != line {
				fmt.Fprintln(app.Output, *app.PreviousLine)
			}
			app.PreviousLine = &line
			duplicates = 0
		}
		duplicates++
		if app.PreviousLine != nil {
			fmt.Fprintln(app.Output, *app.PreviousLine)
		}
	}
}
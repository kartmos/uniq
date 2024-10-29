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
var fFlag int
var sFlag int

type App struct {
	Input              *os.File
	Output             *os.File
	PreviousLine       *string
	PreviousLineOrigin *string
	PrintLine          *string
}

var app = App{}

func init() {
	flag.BoolVar(&cFlag, "c", false, "Подсчитать количество встречаний строки во входных данных")
	flag.BoolVar(&uFlag, "u", false, "Не повторяющиеся строки")
	flag.BoolVar(&dFlag, "d", false, "Повторяющиеся строки")
	flag.BoolVar(&iFlag, "i", false, "Сравнение без учета регистра")
	flag.IntVar(&fFlag, "f", 0, "Пропускаем поле")
	flag.IntVar(&sFlag, "s", 0, "Пропускаем символ")
}

// Выполняет функционал первых 3х флагов если также есть флаг "-i"
func withIArg(str, strOrig string, r int) int {
	var duplicates = r
	switch true {
	case cFlag:
		if app.PreviousLine == nil {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = &str
			duplicates++
		} else if *app.PreviousLine == str && duplicates == 1 {
			app.PrintLine = app.PreviousLineOrigin
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine == str && duplicates > 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine != str {
			fmt.Fprintf(app.Output, "%d %s\n", duplicates, *app.PrintLine)
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		}
	case uFlag:
		if app.PreviousLine == nil {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = &str
			duplicates++
		} else if *app.PreviousLine == str && duplicates == 1 {
			app.PrintLine = app.PreviousLineOrigin
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine == str && duplicates > 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine != str && duplicates > 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		} else if *app.PreviousLine != str && duplicates == 1 {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		}
	case dFlag:
		if app.PreviousLine == nil {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = &str
			duplicates++
		} else if *app.PreviousLine == str && duplicates == 1 {
			app.PrintLine = app.PreviousLineOrigin
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine == str && duplicates > 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine != str && duplicates > 1 {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		} else if *app.PreviousLine != str && duplicates == 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		}
	default:
		if app.PreviousLine == nil {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = &str
			duplicates++
		} else if *app.PreviousLine == str && duplicates == 1 {
			app.PrintLine = app.PreviousLineOrigin
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine == str && duplicates > 1 {
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			duplicates++
		} else if *app.PreviousLine != str {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		}
	}
	return duplicates
}

// Выполняет функционал первых 3х флагов без флага "-i"

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

	// func fArg(inLine, string) string{
	// 	fields := strings.Fields(inline)
	// 	countFields := len(fields)
	// 	for i = *fFlag - 1; i <= countFields; i++
	// }

	stage := bufio.NewScanner(app.Input)
	count := 0
	for stage.Scan() {
		if iFlag {
			line := strings.ToLower(stage.Text())
			lineOrigin := stage.Text()
			count = withIArg(line, lineOrigin, count)
		} else {
			lineOrigin := stage.Text()
			count = withIArg(lineOrigin, lineOrigin, count)
		}
	}
	//Проверка последнего предложения если оно не вывелось
	switch true {
	case cFlag:
		if app.PreviousLineOrigin != nil && count > 0 {
			fmt.Fprintf(app.Output, "%d %s\n", count, *app.PrintLine)
		}
	case uFlag:
		if app.PreviousLineOrigin != nil && count == 1 {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
		}
	case dFlag:
		if app.PreviousLineOrigin != nil && count > 1 {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
		}
	default:
		if app.PreviousLineOrigin != nil && count > 0 {
			fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
		}
	}
}

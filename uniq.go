package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var uniqueFlag bool
var repeatedFlag bool
var countFlag bool
var ignoreCaseFlag bool
var skipFieldsFlag int
var skipCharFlag int
var sumBoolFlag int

type App struct {
	Input              *os.File
	Output             *os.File
	PreviousLine       *string
	PreviousLineOrigin *string
	Count              int
	Scanner            *bufio.Scanner
}

var app = App{}

func init() {
	flag.BoolVar(&countFlag, "c", false, "Подсчитать количество встречаний строки во входных данных")
	flag.BoolVar(&uniqueFlag, "u", false, "Не повторяющиеся строки")
	flag.BoolVar(&repeatedFlag, "d", false, "Повторяющиеся строки")
	flag.BoolVar(&ignoreCaseFlag, "i", false, "Сравнение без учета регистра")
	flag.IntVar(&skipFieldsFlag, "f", 0, "Пропускаем поле")
	flag.IntVar(&skipCharFlag, "s", 0, "Пропускаем символ")
}

func openFile(filename string) *os.File {
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	return fh
}

func skipFields(inPutString string) string {
	indexOfSpace := -1
	numSpace := 0
	for i, char := range inPutString {
		if numSpace == skipFieldsFlag {
			break
		} else if char == ' ' {
			indexOfSpace = i
			numSpace++
		}
	}
	outPutString := inPutString[indexOfSpace+1:]
	return outPutString
}

func skipChar(inPutString string) string {
	if skipCharFlag >= len(inPutString) {
		return ""
	}
	return inPutString[skipCharFlag:]
}

func (r *App) reader() (string, bool) {
	if r.Scanner.Scan() {
		return r.Scanner.Text(), true
	}
	return "", false
}

func (p *App) processer(str, strOrig string) (string, int, bool) {
	if ignoreCaseFlag {
		str = strings.ToLower(skipChar(skipFields(str)))
	} else {
		str = skipChar(skipFields(str))
	}

	if p.PreviousLine == nil {
		p.PreviousLine = &str
		p.PreviousLineOrigin = &strOrig
		p.Count = 1
		return "", 0, false
	}

	if *p.PreviousLine == str {
		p.Count++
		return "", 0, false
	} else {
		lineToPrint := *p.PreviousLineOrigin
		countToPrint := p.Count
		p.PreviousLine = &str
		p.PreviousLineOrigin = &strOrig
		p.Count = 1
		return lineToPrint, countToPrint, true
	}
}

func (w *App) writer(str string, count int) {
	if uniqueFlag && count > 1 {
		return
	}
	if repeatedFlag && count == 1 {
		return
	}

	if countFlag {
		fmt.Fprintf(w.Output, "   %d %s\n", count, str)
	} else {
		fmt.Fprintf(w.Output, "%s\n", str)
	}
}

func (a *App) Run() {
	a.Scanner = bufio.NewScanner(a.Input)
	is_run := true
	for is_run {
		strOrig, run := a.reader()
		if !run {
			break
		}
		strWrite, countWrite, result := a.processer(strOrig, strOrig)
		if result {
			a.writer(strWrite, countWrite)
		}
	}
	if app.Count > 0 {
		app.writer(*app.PreviousLineOrigin, app.Count)
	}
}

func main() {
	flag.Parse()
	if countFlag {
		sumBoolFlag++
	}
	if uniqueFlag {
		sumBoolFlag++
	}
	if repeatedFlag {
		sumBoolFlag++
	}
	if sumBoolFlag > 1 {
		fmt.Fprint(os.Stderr, "Can't use options -c | -d | -u together")
		os.Exit(1)
	}

	if len(flag.Args()) > 0 {
		app.Input = openFile(flag.Args()[0])
		defer app.Input.Close()
	} else {
		app.Input = os.Stdin
	}

	if len(flag.Args()) > 1 {
		app.Output = openFile(flag.Args()[1])
		defer app.Output.Close()
	} else {
		app.Output = os.Stdout
	}
	app.Run()
}

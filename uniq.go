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

// Выполняет функционал первых 3х флагов без флага "-i"

func openFile(filename string) *os.File {
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	return fh
}

func skipFields(inPutString string) string{
	inPutFields := strings.Fields(inPutString)
	outPutFields := []string {}
	for i, field := range inPutFields{
		if i >= fFlag {
			outPutFields = append(outPutFields, field)
		}
	}
	outPutString := strings.Join(outPutFields, " ")
	return outPutString
}

func skipChar(inPutString string) string {
	outPutCharRune := []byte {}
	for i, char := range inPutString {
		if i >= sFlag {
			outPutCharRune = append(outPutCharRune, byte(char))
		}
	}
	outPutString := string(outPutCharRune)
	return outPutString
}

func changeRegistrLow(inPutString string) string {
	return strings.ToLower(inPutString)
}

func sortString(str, strOrig string, r int) int {
	var duplicates = r
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
		} else if *app.PreviousLine != str && (!dFlag && !cFlag  && !uFlag) || cFlag {
				if cFlag {
					fmt.Fprintf(app.Output, "   %d %s\n", duplicates, *app.PrintLine)
				}
				if !dFlag && !cFlag && !uFlag {
					fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
				}
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		} else if *app.PreviousLine != str && duplicates > 1 && (uFlag || dFlag){
			if dFlag {
				fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)
			}
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		} else if *app.PreviousLine != str && duplicates == 1 && (uFlag || dFlag){
			if uFlag{fmt.Fprintf(app.Output, "%s\n", *app.PrintLine)}
			app.PreviousLine = &str
			app.PreviousLineOrigin = &strOrig
			app.PrintLine = app.PreviousLineOrigin
			duplicates = 1
		}
	return duplicates
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

	scanner := bufio.NewScanner(app.Input)
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "%s", scanner.Err())
	}
	count := 0
	// fmt.Println(fFlag)
	for scanner.Scan() {
		originString := scanner.Text()
		if iFlag {
			outFuncString := changeRegistrLow(skipChar(skipFields(scanner.Text())))	
			count = sortString(outFuncString, originString, count)
		} else {
			outFuncString := skipChar(skipFields(scanner.Text()))
			count = sortString(outFuncString, originString, count)
		}
	}

	//Проверка последнего предложения если оно не вывелось

	
	switch true {
	case cFlag:
		if app.PreviousLineOrigin != nil && count > 0 {
			fmt.Fprintf(app.Output, "   %d %s\n", count, *app.PrintLine)
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

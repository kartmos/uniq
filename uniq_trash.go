package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var countLines bool
var sFlag int

func init() {
	flag.BoolVar(&countLines, "c", false, "count")
	flag.IntVar(&sFlag, "s", 0, "skip")
}

type App struct {
	Input        *os.File
	Output       *os.File
	PreviousLine string
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
	fmt.Printf("countLines: %v :%+v\n", countLines, flag.Args())

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

	fmt.Println("Начало обработки файла") // Отладочное сообщение

	stage := bufio.NewScanner(app.Input)
	duplicates := 0

	for stage.Scan() {
		line := stage.Text()
		fmt.Println("Считана строка:", line) // Отладочное сообщение

		if app.PreviousLine == "" {
			app.PreviousLine = line
			duplicates = 1
		} else if app.PreviousLine != line {
			if countLines {
				fmt.Fprintf(app.Output, "%d %s\n", duplicates, app.PreviousLine)
			} else {
				fmt.Fprintln(app.Output, app.PreviousLine)
			}
			app.PreviousLine = line
			duplicates = 1
		} else {
			duplicates++
		}
	}

	// Финальный вывод последней строки после завершения цикла
	if app.PreviousLine != "" {
		if countLines {
			fmt.Fprintf(app.Output, "%d %s\n", duplicates, app.PreviousLine)
		} else {
			fmt.Fprintln(app.Output, app.PreviousLine)
		}
	}

	fmt.Println("Обработка завершена") // Отладочное сообщение
}

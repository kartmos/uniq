file, err := os.Open("text.txt")
if err != nil {
	fmt.Println("No file", err)
	return
}
defer file.Close()

stage := []Sent{}

sc := bufio.NewScanner(file)
for sc.Scan() {
	st = sc.Text()

	fmt.Println(stage)
}


func init() {
    //Задаем правила разбора:
    flag.IntVar(&WORKERS, "w", WORKERS, "количество потоков")
    flag.IntVar(&REPORT_PERIOD, "r", REPORT_PERIOD, "частота отчетов (сек)")
    flag.IntVar(&DUP_TO_STOP, "d", DUP_TO_STOP, "кол-во дубликатов для остановки")
    flag.StringVar(&HASH_FILE, "hf", HASH_FILE, "файл хешей")
    flag.StringVar(&QUOTES_FILE, "qf", QUOTES_FILE, "файл записей")
    //И запускаем разбор аргументов
    flag.Parse() 
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type textFile struct {
	Sent string
	Num  int
}

var countLines bool
var sFlag int

func init() {
	flag.BoolVar(&countLines, "c", false, "count")
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
	fmt.Printf("countLines: %v :%+#v\n", countLines, flag.Args())

	app := App{}

	argN := len(flag.Args())
	switch argN {
	case 0:
		app.Input = os.Stdin
		app.Output = os.Stdout
	case 1:
		app.Input = openFile(flag.Args()[0])
		defer app.Input.Close()
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

	//Создаю слайс куда засуну строки
	slice := []textFile{}
	var n int

	stage := bufio.NewScanner(app.Input)
	duplicates := 0
	//Добавляю элементы в слайс
	for stage.Scan() {
		line := stage.Text()
		if app.PreviousLine == nil {
			app.PreviousLine = &line
		} else if *app.PreviousLine != line {
			if countLines {
				fmt.Fprintf(app.Output, "%d %s\n", duplicates, app.PreviousLine)
			} else {
				fmt.Fprintln(app.Output, app.PreviousLine)
			}
			app.PreviousLine = &line
			duplicates = 0
		}
		duplicates++

		newSlice := textFile{Sent: stage.Text()}
		slice = append(slice, newSlice)
	}

	//Хочу задать повторы предложений цифрами в слайсе струры значению textFile.Num каждому предложению
	for i := range slice {
		if i > 0 && slice[i].Sent != slice[i-1].Sent {
			n = 1
		} else {
			n++
		}
		slice[i].Num = n
	}

	//for _, num := range slice {fmt.Printf("%d %s\n", num.Num, num.Sent)}

	//Задал l и с величину и объем слайса, чтобы потом ипользовать в циклах и сделать программу более гибкой и она работала для любых комбинаций предложений, а также с количеством предложений, но тут есть недочет скорее всего с capacity, так как слайс добавляяет предложение поштучно, то ем больше предложений будет тем больше capacity для слайса GO будет выделять cap X2 от предыдущего значения, что может привести к слишком большому количеству выделянной памяти, но не используемой.

	l := len(slice)
	c := cap(slice)
	if len(os.Args) > 1 {
		fmt.Printf("text.txt | go run main.go %s\n-----------\n", os.Args[1])
		arg := os.Args[1]
		switch arg {
		// Не понял зачем эти параметры вообще нужны но подогнал под результат "input_file" "input_file output_file"
		case "input_file":
			for i := range slice {
				if i < l-1 && slice[i].Sent != slice[i+1].Sent {
					fmt.Printf("%s\n", slice[i].Sent)
				} else if i == l-1 && slice[i].Sent == slice[i-1].Sent {
					fmt.Printf("%s\n", slice[i].Sent)
				}
			}
		case "input_file output_file":
			for i := range slice {
				if i < l-1 && slice[i].Sent != slice[i+1].Sent {
					fmt.Printf("%s\n", slice[i].Sent)
				} else if i == l-1 && slice[i].Sent == slice[i-1].Sent {
					fmt.Printf("%s\n", slice[i].Sent)
				}
			}
		case "-c":
			for i := range slice {
				if i < l-1 && slice[i].Sent != slice[i+1].Sent {
					fmt.Printf("%d %s\n", slice[i].Num, slice[i].Sent)
				} else if i == l-1 && slice[i].Sent == slice[i-1].Sent {
					fmt.Printf("%d %s\n", slice[i].Num, slice[i].Sent)
				}
			}
		case "-d":
			for i := range slice {
				if slice[i].Num == 2 {
					fmt.Printf("%s\n", slice[i].Sent)
				}
			}
		case "-u":
			for i := range slice {
				if i < l-1 && slice[i].Num == 1 && slice[i+1].Num != 2 {
					fmt.Printf("%s\n", slice[i].Sent)
				} else if i == l-1 && slice[i].Num <= slice[i-1].Num {
					fmt.Printf("%s\n", slice[i].Sent)
				}
			}
		case "-i":
			for i := range slice {
				if i < l-1 && strings.ToLower(slice[i].Sent) != strings.ToLower(slice[i+1].Sent) {
					fmt.Printf("%s\n", slice[i].Sent)
				} else if i == l-1 && strings.ToLower(slice[i].Sent) == strings.ToLower(slice[i-1].Sent) {
					fmt.Printf("%s\n", slice[i].Sent)
				}
			}
			// case "-f":
			// 	fmt.Printf("\nАргумуент %s Num %d\n\n", arg, numField)
			// 	for i := range slice {
			// 		if i < l-1 {
			// 			a := strings.Fields(slice[i].Sent)
			// 			b := strings.Fields(slice[i+1].Sent)
			// 			// lenA :=
			// 			if compareFields(a, b) == true {
			// 				fmt.Printf("%s\n", slice[i].Sent)
			// 			}
			// 		} else if i == l-1 {
			// 			a := strings.Fields(slice[i].Sent)
			// 			b := strings.Fields(slice[i-1].Sent)
			// 			if compareFields(a, b) == true {
			// 				fmt.Printf("%s\n", slice[i].Sent)
			// 			}
			// 		}
			// 	}
		}

	} else {
		for i := range slice {
			if i < l-1 && slice[i].Sent != slice[i+1].Sent {
				fmt.Printf("%s\n", slice[i].Sent)
			} else if i == l-1 && slice[i].Sent == slice[i-1].Sent {
				fmt.Printf("%s\n", slice[i].Sent)
			}
		}
		fmt.Printf("-----------\nCapacity %d\nLenth %d\n\n", c, l)
	}
}

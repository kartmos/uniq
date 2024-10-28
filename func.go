func firstArg (scanner *bufio.Scanner) {

    switch true {
    case cFlag:
		fmt.Printf("\ncFlag: %v :%+v\n", cFlag, flag.Args())
		fmt.Printf("\nПосчитать строки\n\n")
		for stage.Scan() {
			line := stage.Text()
			if app.PreviousLine == nil{
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

}

func iArg
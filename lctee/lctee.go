package lctee

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func Filter(filePointer *os.File, parser LogcatParser, printer LogcatPrinter) {
	scanner := bufio.NewScanner(filePointer)

	for scanner.Scan() {
		line := scanner.Text()
		formatPointer := parser.parse(line)

		if formatPointer != nil {
			printer.print(formatPointer)
		} else {
			fmt.Printf("error: %s\n", line)
		}
	}
}

func FilterWithFile(file string, parser LogcatParser, printer LogcatPrinter) {
	var filePointer *os.File
	var err error

	filePointer, err = os.Open(file)
	if nil != err {
		panic(err)
	}
	defer filePointer.Close()

	Filter(filePointer, parser, printer)
}

func FilterWithReader(readCloser io.ReadCloser, parser LogcatParser, printer LogcatPrinter) {
	reader := bufio.NewReader(readCloser)
	for {
		line, _, err := reader.ReadLine()
		if nil != err {
			break
		}

		formatPointer := parser.parse(string(line))
		if nil != formatPointer {
			printer.print(formatPointer)
		} else {
			fmt.Printf("error: %s\n", line)
		}
	}
}

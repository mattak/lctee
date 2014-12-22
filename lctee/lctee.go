package lctee

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"io"
)

const (
	StringColorVerbose = 35
	StringColorDebug   = 34
	StringColorInfo    = 32
	StringColorWarn    = 33
	StringColorError   = 31
)

type LogcatParser interface {
	parse(string) *LogcatFormat
}

type LogcatPrinter interface {
	print(*LogcatFormat)
}

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

type LogcatParserThread struct {
	Description string
}

func (parser LogcatParserThread) parse(line string) *LogcatFormat {
	// unformated message. such as: "--------- beginning of crash"
	if len(line) < 33 {
		return &LogcatFormat{"00-00 00:00:00.000", "00000", "00000", "V", "#", line}
	}

	var time = line[0:18]
	var processId = line[19:24]
	var threadId = line[25:30]
	var level = line[31:32]
	var tagAndMessage = line[33:]

	tagAndMessage = strings.Replace(tagAndMessage, "\t", " ", -1)
	separaterIndex := strings.Index(tagAndMessage, ":")

	// separater doesnot exists
	if -1 == separaterIndex {
		return &LogcatFormat{"00-00 00:00:00.000", "00000", "00000", "V", "#", line}
	}

	tag := tagAndMessage[:separaterIndex]
	message := tagAndMessage[(separaterIndex + 2):]

	return &LogcatFormat{time, processId, threadId, level, tag, message}
}

func RemoveSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func CreatePaddedTag(tag string) string {
	var tagPadded string

	// There exists over 23 characters tag. such as com.amazon.kindle/com.amazon.identity.auth.device.utils.CentralApkUtils
	if len(tag) > 23 {
		tagPadded = tag[:24]
	} else {
		tagPaddingLength := 23 - len(tag)

		if tagPaddingLength > 0 {
			tagPadded = tag + strings.Repeat(" ", tagPaddingLength)
		} else {
			tagPadded = tag
		}
	}

	return tagPadded
}

func PrettyLogcatFormat(format *LogcatFormat) *LogcatFormat {
	format.ProcessId = RemoveSpace(format.ProcessId)
	format.ThreadId = RemoveSpace(format.ThreadId)
	format.Tag = RemoveSpace(format.Tag)
	return format
}

func ColoringString(color int, body string) string {
	return fmt.Sprintf("\033[%dm%s\033[m", color, body)
}

func ColoringStringWithFormat(format *LogcatFormat, message string) string {
	switch format.Level {
	case "V":
		message = ColoringString(StringColorVerbose, message)
	case "D":
		message = ColoringString(StringColorDebug, message)
	case "I":
		message = ColoringString(StringColorInfo, message)
	case "W":
		message = ColoringString(StringColorWarn, message)
	case "E":
		message = ColoringString(StringColorError, message)
	}
	return message
}

type LogcatPrinterDefault struct {
	Color          bool
	PreviousFormat LogcatFormat
}

type LogcatPrinterLTSV struct {
	Color bool
}

type LogcatPrinterJSON struct {
	Color bool
}

func (filter LogcatPrinterLTSV) print(format *LogcatFormat) {
	var message = fmt.Sprintf("time:%s\tpid:%s\ttid:%s\tlevel:%s\ttag:%s\tmessage:%s\n",
		format.Time, format.ProcessId, format.ThreadId, format.Level, format.Tag, format.Message)

	if filter.Color {
		message = ColoringStringWithFormat(format, message)
	}

	fmt.Println(message)
}

func (filter LogcatPrinterJSON) print(format *LogcatFormat) {
	PrettyLogcatFormat(format)
	data, _ := json.Marshal(*format)
	message := string(data)

	if filter.Color {
		message = ColoringStringWithFormat(format, message)
	}

	fmt.Println(message)
}

func (filter LogcatPrinterDefault) print(format *LogcatFormat) {
	PrettyLogcatFormat(format)
	tag := CreatePaddedTag(format.Tag)
	var message string

	if "" == filter.PreviousFormat.Tag || format.Tag != filter.PreviousFormat.Tag {
		message = fmt.Sprintf("%s \t%s", tag, format.Message)
	} else {
		message = fmt.Sprintf("                        \t%s", format.Message)
	}

	if filter.Color {
		message = ColoringStringWithFormat(format, message)
	}

	fmt.Println(message)
	filter.PreviousFormat.Tag = format.Tag
}

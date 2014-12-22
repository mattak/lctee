package lctee

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LogcatPrinter interface {
	print(*LogcatFormat)
}

func removeSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func createPaddedTag(tag string) string {
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

func prettyLogcatFormat(format *LogcatFormat) *LogcatFormat {
	format.ProcessId = removeSpace(format.ProcessId)
	format.ThreadId = removeSpace(format.ThreadId)
	format.Tag = removeSpace(format.Tag)
	return format
}

func coloringString(color int, body string) string {
	return fmt.Sprintf("\033[%dm%s\033[m", color, body)
}

func coloringStringWithFormat(format *LogcatFormat, message string) string {
	switch format.Level {
	case "V":
		message = coloringString(StringColorVerbose, message)
	case "D":
		message = coloringString(StringColorDebug, message)
	case "I":
		message = coloringString(StringColorInfo, message)
	case "W":
		message = coloringString(StringColorWarn, message)
	case "E":
		message = coloringString(StringColorError, message)
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
		message = coloringStringWithFormat(format, message)
	}

	fmt.Println(message)
}

func (filter LogcatPrinterJSON) print(format *LogcatFormat) {
	prettyLogcatFormat(format)
	data, _ := json.Marshal(*format)
	message := string(data)

	if filter.Color {
		message = coloringStringWithFormat(format, message)
	}

	fmt.Println(message)
}

func (filter LogcatPrinterDefault) print(format *LogcatFormat) {
	prettyLogcatFormat(format)
	tag := createPaddedTag(format.Tag)
	var message string

	if "" == filter.PreviousFormat.Tag || format.Tag != filter.PreviousFormat.Tag {
		message = fmt.Sprintf("%s \t%s", tag, format.Message)
	} else {
		message = fmt.Sprintf("                        \t%s", format.Message)
	}

	if filter.Color {
		message = coloringStringWithFormat(format, message)
	}

	fmt.Println(message)
	filter.PreviousFormat.Tag = format.Tag
}

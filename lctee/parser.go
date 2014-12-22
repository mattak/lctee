package lctee

import (
  "strings"
)

type LogcatParser interface {
	parse(string) *LogcatFormat
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

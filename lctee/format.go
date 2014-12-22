package lctee

type LogcatFormat struct {
	Time      string `json:"time"`
	ProcessId string `json:"pid"`
	ThreadId  string `json:"tid"`
	Level     string `json:"level"`
	Tag       string `json:"tag"`
	Message   string `json:"message"`
}

type LogcatQueue struct {
	list []LogcatFormat
}

func (queue *LogcatQueue) Push(format *LogcatFormat) {
}

func (queue *LogcatQueue) Remove() *LogcatFormat {
	return nil
}

func (queue *LogcatQueue) Fetch(pointBefore int) *LogcatFormat {
	return nil
}

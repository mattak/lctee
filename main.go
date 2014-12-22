package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/mattak/lctee/lctee"
	"os"
	"os/exec"
	"path/filepath"
)

func scan(filelist []string, parser lctee.LogcatParser, printer lctee.LogcatPrinter) {
	if len(filelist) <= 0 {
		lctee.Filter(os.Stdin, parser, printer)
	} else {
		for _, file := range filelist {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Not exist file: %s\n", file)
				os.Exit(1)
			} else {
				lctee.FilterWithFile(file, parser, printer)
			}
		}
	}
}

func getLogHome() string {
	loghome := os.Getenv("LCTEE_LOGHOME")
	if len(loghome) > 0 {
		return loghome
	}

	home := os.Getenv("HOME")
	return home + "/.lctee/logs"
}

func isFile(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}

	defer file.Close()

	filestat, err := file.Stat()
	if err != nil {
		return false
	}
	switch mode := filestat.Mode(); {
	case mode.IsDir():
		return false
	case mode.IsRegular():
		return true
	default:
		return false
	}
}

func visitLogFile(path string, f os.FileInfo, err error) error {
	if isFile(path) {
		println(path)
	}
	return nil
}

func TaskDefault(c *cli.Context) {
	if c.Bool("loghome") {
		println(getLogHome())
		os.Exit(0)
	}
	withColor := !c.Bool("no-color")

	cmd := exec.Command("adb", "logcat", "-v", "threadtime")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		println(err)
		os.Exit(1)
	}

	cmd.Start()

	parser := lctee.LogcatParserThread{""}
	printer := lctee.LogcatPrinterDefault{withColor, lctee.LogcatFormat{}}

	lctee.FilterWithReader(stdout, parser, printer)
}

func TaskClear(c *cli.Context) {
	cmd := exec.Command("adb", "logcat", "-c")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		println(err)
		os.Exit(1)
	}

	println("cleared")
}

func TaskFile(c *cli.Context) {
	loghome := getLogHome()
	if _, err := os.Stat(loghome); err == nil {
		filepath.Walk(loghome, visitLogFile)
	}
}

func TaskInput(c *cli.Context) {
	parser := lctee.LogcatParserThread{""}
	withColor := !c.Bool("no-color")
	filelist := c.Args()

	switch c.String("format") {
	case "ltsv":
		printer := lctee.LogcatPrinterLTSV{withColor}
		scan(filelist, parser, printer)
	case "json":
		printer := lctee.LogcatPrinterJSON{withColor}
		scan(filelist, parser, printer)
	case "pidcat":
		fallthrough
	default:
		printer := lctee.LogcatPrinterDefault{withColor, lctee.LogcatFormat{}}
		scan(filelist, parser, printer)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "lctee"
	app.Usage = "logcat and run tee."
	app.Version = "0.0.1"
	app.Author = "mattak"
	app.Email = "mattak.me@gmail.com"
	app.Action = TaskDefault
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "loghome",
			Usage: "Show log home.",
		},
		cli.BoolFlag{
			Name:  "no-color, n",
			Usage: "Print with color.",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "input",
			ShortName: "i",
			Usage:     "Input external logcat data from stdin or files.",
			Action:    TaskInput,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format, f",
					Value: "pidcat",
					Usage: "Set output format, pidcat,json,ltsv.",
				},
				cli.BoolFlag{
					Name:  "no-color, n",
					Usage: "Print with color.",
				},
			},
		},
		{
			Name:      "clear",
			ShortName: "c",
			Usage:     "Clear the logs. adb logcat -c",
			Action:    TaskClear,
		},
		{
			Name:      "file",
			ShortName: "f",
			Usage:     "Show log files",
			Action:    TaskFile,
		},
	}
	app.Run(os.Args)
}

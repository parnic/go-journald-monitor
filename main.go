package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var (
	unit = flag.String("unit", "", "specifies a unit to filter the journald output for")
)

const (
	cursorPrefix       = "-- cursor: "
	cursorFilenameBase = "lastCursor"
)

func main() {
	flag.Parse()

	parseJournalData(getJournalData())
}

func getJournalData() (cursorFilename, lastCursor string, output []byte) {
	cursorFilename, lastCursor = getLastCursor()

	var err error
	if output, err = getCmd(lastCursor).Output(); err != nil {
		panic(err)
	}

	return
}

func getLastCursor() (cursorFilename, lastCursor string) {
	cursorFilename = cursorFilenameBase
	if len(*unit) > 0 {
		cursorFilename = fmt.Sprintf("%s-%s", cursorFilenameBase, *unit)
	}

	lastCursorBytes, _ := ioutil.ReadFile(cursorFilename)
	lastCursor = string(lastCursorBytes)

	return
}

func getCmd(lastCursor string) (cmd *exec.Cmd) {
	cmd = exec.Command("journalctl", "--quiet", "--show-cursor", "--no-pager")
	if len(*unit) > 0 {
		cmd.Args = append(cmd.Args, []string{
			"-u",
			*unit,
		}...)
	}
	if len(lastCursor) > 0 {
		cmd.Args = append(cmd.Args, fmt.Sprintf("--after-cursor=%s", lastCursor))
	}

	return
}

func parseJournalData(cursorFilename, lastCursor string, output []byte) {
	buf := bytes.NewBuffer(output)
	scanner := bufio.NewScanner(buf)

	for hasData := true; hasData; hasData = scanner.Scan() {
		handleLine(scanner.Text(), lastCursor, cursorFilename)
	}
}

func handleLine(line, lastCursor, cursorFilename string) {
	if strings.HasPrefix(line, cursorPrefix) {
		lastCursor = line[len(cursorPrefix):]
		ioutil.WriteFile(cursorFilename, []byte(lastCursor), 0664)
	} else if len(lastCursor) > 0 && len(line) > 0 {
		fmt.Println(line)
	}
}

package mooze

import (
	"io/ioutil"
	"os"
)

type HistoryWriter struct {
	file *os.File
}

func NewHistoryWriter() *HistoryWriter {
	return &HistoryWriter{}
}

func (h *HistoryWriter) Write(s string) {
	err := ioutil.WriteFile("history", []byte(s), 0644)
	if err != nil {
		panic(err)
	}
}

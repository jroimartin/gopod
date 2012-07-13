package podcast

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type PodcastLog struct {
	file string
}

func NewPodcastLog(file string) *PodcastLog {
	return &PodcastLog{file: file}
}

func (l *PodcastLog) CheckLog(url string) (bool, error) {
	f, err := os.Open(l.file)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		if strings.Contains(string(line), url) {
			return true, nil
		}
	}
	return false, nil
}

func (l *PodcastLog) AddLog(url string)  error {
	flags := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	f, err := os.OpenFile(l.file, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, url)
	if err != nil {
		return err
	}
	return nil
}

package podcast

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type PodcastList struct {
	file string
}

func NewPodcastList(file string) *PodcastList {
	return &PodcastList{file: file}
}

func (l *PodcastList) Get() ([]string, error) {
	f, err := os.Open(l.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("No podcasts")
		}
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	var podcasts []string
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			return podcasts, nil
		}
		if err != nil {
			return nil, err
		}
		podcasts = append(podcasts, string(line))
	}
	return podcasts, nil
}

func (l *PodcastList) Add(rss string) error {
	flags := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	f, err := os.OpenFile(l.file, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, rss)
	if err != nil {
		return err
	}
	return nil
}

func (l *PodcastList) Remove(n int) error {
	podcasts, err := l.Get()
	if err != nil {
		return err
	}
	f, err := os.Create(l.file)
	if err != nil {
		return err
	}
	defer f.Close()
	for i, p := range podcasts {
		if i == n {
			continue
		}
		_, err = fmt.Fprintln(f, p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *PodcastList) String() string {
	var s string
	podcasts, err := l.Get()
	if err != nil {
		return ""
	}
	for i, p := range podcasts {
		s += fmt.Sprintf("[%d] %s\n", i, p)
	}
	return s
}

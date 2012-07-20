package podcast

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type List struct {
	File string
	Items []string
}

func NewList() *List {
	return &List{}
}

func Open(file string) (*List, error) {
	l := NewList()
	l.File = file
	l.Items = make([]string, 0)
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return l, nil
		}
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			return l, nil
		}
		if err != nil {
			return nil, err
		}
		l.Items = append(l.Items, string(line))
	}
	return l, nil
}

func (l *List) Dump() error {
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	f, err := os.OpenFile(l.File, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, item := range l.Items {
		_, err = fmt.Fprintln(f, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *List) Add(url string) error {
	if l.Exists(url) {
		return errors.New("duplicated podcast")
	}
	l.Items = append(l.Items, url)
	return nil
}

func (l *List) Remove(n int) error {
	if n < 0 || n >= len(l.Items) {
		return errors.New("index out of bounds")
	}
	tmp := make([]string, 0)
	for i, item := range l.Items {
		if i == n {
			continue
		}
		tmp = append(tmp, item)
	}
	l.Items = tmp
	return nil
}

func (l *List) Exists(url string) bool {
	for _, item := range l.Items {
		if strings.Contains(item, url) {
			return true
		}
	}
	return false
}

func (l *List) String() string {
	var s string
	for i, item := range l.Items {
		s += fmt.Sprintf("[%d] %s\n", i, item)
	}
	return s
}

package podcast

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var PrintStatus func(written, total int64) = nil

type Enclosure struct {
	Url  string `xml:"url,attr"`
	Type string `xml:"type,attr"`
}

type Episode struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Enclosure   Enclosure `xml:"enclosure"`
	PubDate     string    `xml:"pubDate"`
	Duration    string    `xml:"duration"`
}

type PodcastXML struct {
	XMLName     xml.Name  `xml:"rss"`
	Title       string    `xml:"channel>title"`
	Link        string    `xml:"channel>link"`
	Description string    `xml:"channel>description"`
	Episodes    []Episode `xml:"channel>item"`
}

type Podcast struct {
	XML PodcastXML
	Url string
}

func NewPodcast(url string) *Podcast {
	return &Podcast{Url: url}
}

func (p *Podcast) Get() error {
	resp, err := http.Get(p.Url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d := xml.NewDecoder(resp.Body)
	err = d.Decode(&p.XML)
	if err != nil {
		return err
	}
	return nil
}

func wget(dst io.Writer, src io.Reader, total int64) error {
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				return ew
			}
			if nr != nw {
				return errors.New("Bytes read != Bytes written")
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			return er
		}
		if PrintStatus != nil {
			PrintStatus(written, total)
		}
	}
	return nil
}

func (e *Episode) Download(folder string) error {
	url := e.Enclosure.Url
	fname := filepath.Join(folder, filepath.Base(url))

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = wget(f, r.Body, r.ContentLength)
	if err != nil {
		return err
	}
	return nil
}

func (p *Podcast) String() string {
	s := fmt.Sprintln("Title:", p.XML.Title)
	s += fmt.Sprintln("Link:", p.XML.Link)
	s += fmt.Sprintln("Description:", p.XML.Description)
	s += fmt.Sprintln("Episodes:", len(p.XML.Episodes))
	return s
}

package podcast

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

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

	_, err = io.Copy(f, r.Body)
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

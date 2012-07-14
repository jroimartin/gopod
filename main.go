package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gopod/podcast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	defaultFolder = path.Join(os.Getenv("HOME"), "podcasts")
	defaultConfig = filepath.Join(os.Getenv("HOME"), ".gopodrc")
	defaultLog    = filepath.Join(os.Getenv("HOME"), ".gopod_log")
	folder        = flag.String("folder", defaultFolder, "folder to store podcasts")
	config        = flag.String("config", defaultConfig, "file to store rss list")
	log           = flag.String("log", defaultLog, "file to track downloaded episodes")
	add           = flag.String("a", "", "add a new podcast")
	remove        = flag.Int("r", -1, "remove a podcast")
	info          = flag.Int("i", -1, "show podcast info")
	list          = flag.Bool("l", false, "list podcasts")
	sync          = flag.Bool("s", false, "sync podcasts")
)

func printStatus(written, total int64) {
	percent := (float64(written) / float64(total)) * 100
	bar := strings.Repeat("=", int(percent/10.0))
	fmt.Fprintf(os.Stderr, "\r%d%% [%-10s] %d/%d", int(percent), bar, written, total)
}

func downloadPodcast(rss string) error {
	p := podcast.NewPodcast(rss)
	err := p.Get()
	if err != nil {
		return err
	}
	l := podcast.NewPodcastList(*log)
	for i, e := range p.XML.Episodes {
		exists, err := l.Check(e.Enclosure.Url)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		fmt.Fprintf(os.Stderr, "Downloading [%d] %s...\n", i, e.Enclosure.Url)
		err = e.Download(*folder)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "\n")
		err = l.Add(e.Enclosure.Url)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncPodcast(l *podcast.PodcastList) error {
	podcasts, err := l.Get()
	if err != nil {
		return err
	}
	err = os.Mkdir(*folder, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	for _, p := range podcasts {
		err = downloadPodcast(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func showInfo(l *podcast.PodcastList, n int) error {
	podcasts, err := l.Get()
	if err != nil {
		return err
	}
	p := podcast.NewPodcast(podcasts[n])
	err = p.Get()
	if err != nil {
		return err
	}
	fmt.Println(p)
	return nil
}

func main() {
	flag.Parse()
	l := podcast.NewPodcastList(*config)
	var err error
	switch {
	case *add != "":
		err = l.Add(*add)
	case *remove != -1:
		err = l.Remove(*remove)
	case *info != -1:
		err = showInfo(l, *info)
	case *sync:
		podcast.PrintStatus = printStatus
		err = syncPodcast(l)
	case *list:
		fmt.Print(l)
	default:
		flag.Usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

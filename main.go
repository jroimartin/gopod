package main

import (
	"errors"
	"flag"
	"log"
	"fmt"
	"github.com/jroimartin/gopod/podcast"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	defaultFolder = path.Join(os.Getenv("HOME"), "podcasts")
	defaultConfigFile = filepath.Join(os.Getenv("HOME"), ".gopodrc")
	defaultLogFile    = filepath.Join(os.Getenv("HOME"), ".gopod_log")
	folder        = flag.String("folder", defaultFolder, "folder to store podcasts")
	configFile        = flag.String("config", defaultConfigFile, "file to store rss list")
	logFile           = flag.String("log", defaultLogFile, "file to track downloaded episodes")
	add           = flag.String("a", "", "add a new podcast")
	remove        = flag.Int("r", -1, "remove a podcast")
	info          = flag.Int("i", -1, "show podcast info")
	list          = flag.Bool("l", false, "list podcasts")
	sync          = flag.Bool("s", false, "sync podcasts")
	all           = flag.Bool("A", false, "mark all podcasts as downloaded")
	quiet         = flag.Bool("q", false, "be quiet while syncing")
)

var lrss, llog *podcast.List

func main() {
	var err error

	flag.Parse()

	lrss, err = podcast.Open(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer lrss.Dump()
	llog, err = podcast.Open(*logFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer llog.Dump()

	switch {
	case *add != "":
		err = lrss.Add(*add)
	case *remove != -1:
		err = lrss.Remove(*remove)
	case *info != -1:
		err = showInfo(*info)
	case *sync:
		if !*quiet {
			podcast.PrintStatus = printStatus
		}
		err = syncAll()
	case *all:
		err = logAll()
	case *list:
		fmt.Print(lrss)
	default:
		flag.Usage()
		os.Exit(2)
	}

	if err != nil {
		log.Fatalln("Error:", err)
	}
}

func showInfo(n int) error {
	if n < 0 || n >= len(lrss.Items) {
		return errors.New("index out of bounds")
	}
	p := podcast.NewPodcast(lrss.Items[n])
	err := p.Get()
	if err != nil {
		return err
	}
	fmt.Println(p)
	return nil
}

func printStatus(written, total int64) {
	percent := (float64(written) / float64(total)) * 100.0
	bar := strings.Repeat("=", int(percent/10.0))
	fmt.Fprintf(os.Stderr, "\r%d%% [%-10s] %d/%d", int(percent), bar, written, total)
}

func syncAll() error {
	err := os.Mkdir(*folder, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}
	for _, rss := range lrss.Items {
		err = syncPodcast(rss)
		if err != nil {
			return err
		}
	}
	return nil
}

func syncPodcast(rss string) error {
	p := podcast.NewPodcast(rss)
	err := p.Get()
	if err != nil {
		return err
	}
	for i, e := range p.XML.Episodes {
		if llog.Exists(e.Enclosure.Url) {
			continue
		}
		if !*quiet {
			fmt.Fprintf(os.Stderr, "Downloading [%d] %s...\n", i, e.Enclosure.Url)
		}
		err = e.Download(*folder)
		if err != nil {
			return err
		}
		if !*quiet {
			fmt.Fprintf(os.Stderr, "\n")
		}
		err = llog.Add(e.Enclosure.Url)
		if err != nil {
			return err
		}
	}
	return nil
}

func logAll() error {
	for _, rss := range lrss.Items {
		err := logPodcast(rss)
		if err != nil {
			return err
		}
	}
	return nil
}

func logPodcast(rss string) error {
	p := podcast.NewPodcast(rss)
	err := p.Get()
	if err != nil {
		return err
	}
	for _, e := range p.XML.Episodes {
		if llog.Exists(e.Enclosure.Url) {
			continue
		}
		err = llog.Add(e.Enclosure.Url)
		if err != nil {
			return err
		}
	}
	return nil
}

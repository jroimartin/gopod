package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gopod/podcast"
	"os"
	"path"
	"path/filepath"
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

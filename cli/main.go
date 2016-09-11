package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/thepaulm/tunecon/sync"
)

type Config struct {
	Src string
	Dst string
}

func getConfig() (Config, error) {
	var cfg Config
	flag.StringVar(&cfg.Src, "src", "", "Source directory")
	flag.StringVar(&cfg.Dst, "dst", "", "Destination directory")
	flag.Parse()
	if cfg.Src == "" {
		return cfg, errors.New("No src directory defined.")
	}
	if cfg.Dst == "" {
		return cfg, errors.New("No dst directory defined.")
	}
	return cfg, nil
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		fmt.Printf("Error parsing config.\n")
		os.Exit(-1)
	}
	sync.Dirs(cfg.Src, cfg.Dst)
}

package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/kyslik/gaglog"
)

func main() {
	filter := &gaglog.GagFilter{
		Writer: os.Stderr,
		Gags: gaglog.Gags{
			regexp.MustCompile("P([a-z]+)ch"): time.Millisecond * 100,
			regexp.MustCompile("L([a-z]+)ch"): time.Millisecond * 1000,
			regexp.MustCompile("F([a-z]+)ch"): time.Millisecond * 10000,
		},
	}

	log.SetOutput(filter)

	for i := 0; i < 1000000; i++ {
		log.Print("Pinch")
		log.Print("Lynch")
		log.Print("Flinch")
		if i%100000 == 0 {
			log.Print("Grinch")
		}
	}
	log.Println("Clinch")
}

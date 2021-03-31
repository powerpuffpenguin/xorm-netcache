package main

import (
	"log"

	"github.com/powerpuffpenguin/xormcache/cmd"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if e := cmd.Execute(); e != nil {
		log.Fatalln(e)
	}
}

package main

import (
	"embed"
	"github.com/0xsj/atelier-wails/root"
	"log"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := root.Run(assets); err != nil {
		log.Fatal(err)
	}
}

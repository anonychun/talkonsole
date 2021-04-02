package main

import (
	"flag"

	"github.com/anonychun/talkonsole/server"
)

func main() {
	port := flag.Int("port", 1401, "")
	flag.Parse()

	err := server.Start(*port)
	if err != nil {
		panic(err.Error())
	}
}

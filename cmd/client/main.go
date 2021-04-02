package main

import (
	"flag"
	"fmt"

	"github.com/anonychun/talkonsole/client"
)

func main() {
	host := flag.String("host", "127.0.0.1", "")
	port := flag.Int("port", 1401, "")
	name := flag.String("name", "", "")
	room := flag.String("room", "public", "")
	flag.Parse()

	if *name == "" {
		fmt.Println("'--help' for help")
		return
	}

	err := client.Join(*host, *port, *name, *room)
	if err != nil {
		panic(err.Error())
	}
}

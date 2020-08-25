package main

import (
	"bufio"
	"collector"
	"log"
	"os"
	"sync"
)

func main() {
	c, _ := collector.NewCollector()
	//c.WatchChannel("UC3IZKseVpdzPSBaWxBxundA")
	//return

	file, err := os.Open("src/collector/channels_on_watch")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var wg sync.WaitGroup
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		wg.Add(1)
		go c.WatchChannel(string(line))
	}

	wg.Wait()
}

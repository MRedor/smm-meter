package main

import (
	"collector"
	"time"
)

func main() {
	c, _ := collector.NewCollector()

	for {
		go c.PopularVideos()
		time.Sleep(6 * time.Hour)
	}
}


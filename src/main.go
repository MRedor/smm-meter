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

	//err := c.ChannelByUsername("dj47maryn")
	//collector.ChannelByUsername(service, "dj47maryn")
	//err := collector.PopularVideos(service)
	//collector.RelatedVideos(service, "")
	//if err != nil {
	//	// TODO
	//}

}


package main

import (
	"bufio"
	"collector"
	"fmt"
	"os"
)

func addChannelsByNames() {
	file, err := os.Open("src/channel_names")
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)

	c, _ := collector.NewCollector()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		c.ChannelByUsername(string(line))
	}


}

func addChannelByIds() {
	file, err := os.Open("src/channel_ids")
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)

	c, _ := collector.NewCollector()
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		c.ChannelById(string(line))
	}
}

func main() {
	addChannelsByNames()
	addChannelByIds()
}

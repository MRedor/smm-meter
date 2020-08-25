package collector

import (
	"db"
	"fmt"
	"time"
)

func (c *Collector) WatchChannel(channelId string) error {
	search := c.service.Search.List([]string{"id,snippet"}).ChannelId(channelId).PublishedAfter(time.Now().AddDate(0, -1, 0).Format(time.RFC3339)).Order("date").MaxResults(50)
	response, err := search.Do()
	if err != nil {
		return err
	}

	var videoIds = ""

	for _, v := range response.Items {
		vidR, err := c.service.Videos.List([]string{VideoParts}).Id(v.Id.VideoId).Do()
		if err != nil {
			return  err
		}
		if len(vidR.Items) == 0 {
			continue
		}
		if !db.VideoExists(vidR.Items[0]) {
			db.AddVideo(vidR.Items[0])
			if videoIds == "" {
				videoIds = vidR.Items[0].Id
			} else {
				videoIds = videoIds + "," + vidR.Items[0].Id
			}
		}
	}

	go c.watchVideoStats(&videoIds)

	for {
		time.Sleep(12 * time.Hour)

		search := c.service.Search.List([]string{"id"}).ChannelId(channelId).PublishedAfter(time.Now().AddDate(0, 0, -1).Format(time.RFC3339)).Order("date").MaxResults(50)
		response, err := search.Do()
		if err != nil {
			return err
		}

		if len(response.Items) > 0 {
			for _, v := range response.Items {
				vidR, err := c.service.Videos.List([]string{VideoParts}).Id(v.Id.VideoId).Do()
				if err != nil {
					return  err
				}
				if !db.VideoExists(vidR.Items[0]) {
					db.AddVideo(vidR.Items[0])
					if videoIds == "" {
						videoIds = vidR.Items[0].Id
					} else {
						videoIds = vidR.Items[0].Id + "," + videoIds
					}
				}
			}
		}
	}

	return nil
}

func (c *Collector) watchVideoStats(videoIds *string) {
	for {
		if *videoIds != "" {
			resp, err := c.service.Videos.List([]string{VideoParts}).Id(*videoIds).Do()
			if err != nil {
				continue
			}
			for _, v := range resp.Items {
				db.AddReactions(v)
			}

			chR, err := c.service.Channels.List([]string{ChannelParts}).Id(resp.Items[0].Snippet.ChannelId).Do()
			if err == nil {
				db.AddChannelStats(chR.Items[0])
				fmt.Println("add stats for channel: " + resp.Items[0].Snippet.ChannelId)
			}
		}

		time.Sleep(60 * time.Minute)
	}
}
package collector

import (
	"db"
	"golang.org/x/exp/errors"
	"time"
)

func (c *Collector) WatchVideo(videoId string, period int) error {

	resp, err := c.service.Videos.List([]string{VideoParts}).Id(videoId).Do()
	if err != nil {
		return err
	}
	if len(resp.Items) == 0 {
		return  errors.New("no such video")
	}
	db.AddVideo(resp.Items[0])
	db.AddReactions(resp.Items[0])

	for {
		time.Sleep(time.Duration(period) * time.Minute)

		resp, err := c.service.Videos.List([]string{VideoParts}).Id(videoId).Do()
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			return errors.New("video was deleted or hidden")
		}
		db.AddReactions(resp.Items[0])
	}
}

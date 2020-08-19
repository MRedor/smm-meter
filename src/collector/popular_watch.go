package collector

import (
	"db"
	"log"
	"os"
	"time"
)

func (c *Collector) PopularVideos() error {
	videos, err := c.service.Videos.List([]string{VideoParts}).Chart("mostPopular").MaxResults(100).Do()
	if err != nil {
		return err
	}
	log.Println("Got new populars: ", len(videos.Items))

	for _, v := range videos.Items {

		if !db.VideoExists(v) {
			err = db.AddVideo(v)
			if err != nil {
				//TODO
			}
		}
		if !db.VideoWasPopular(v) {
			go c.WatchVideo(v.Id)
		}

		db.AddPopular(v)
		//db.AddReactions(v)
	}

	return nil
}

func (c *Collector) WatchVideo(videoId string) {
	f, err := os.OpenFile("src/collector/trends_on_watch", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(videoId + "\n"); err != nil {
		log.Println(err)
	}


	var sleepMinutes int = 15
	var prevErr error
	var views = 0

	for {
		resp, err := c.service.Videos.List([]string{VideoParts}).Id(videoId).Do()
		if err != nil {
			if prevErr != nil {
				break
			} else {
				prevErr = err
				continue
			}
		} else {
			prevErr = nil
		}
		v := resp.Items[0]
		db.AddReactions(v)

		if float64(v.Statistics.ViewCount) > 1.1 * float64(views) {
			sleepMinutes = sleepMinutes / 2
		}
		if float64(v.Statistics.ViewCount) < 1.01 * float64(views) {
			sleepMinutes = sleepMinutes * 2
		}
		if sleepMinutes < 10 {
			if !c.checkIsInPopular(videoId) {
				break
			}
			sleepMinutes = 10
		}
		if sleepMinutes > 60 {
			sleepMinutes = 60
		}

		time.Sleep(time.Duration(sleepMinutes) * time.Minute)
	}
}


func (c *Collector) checkIsInPopular(videoId string) bool {
	videos, _ := c.service.Videos.List([]string{"id"}).Chart("mostPopular").MaxResults(50).Do()

	for _, v := range videos.Items {
		if v.Id == videoId {
			return true
		}
	}

	return false
}
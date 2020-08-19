package collector

import (
	"db"
	"errors"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"time"
)

type Collector struct {
	service *youtube.Service
}

const ChannelParts = "brandingSettings,contentDetails,contentOwnerDetails," +
	"id,invideoPromotion,localizations,snippet,statistics,status,topicDetails"

const VideoParts = "contentDetails,id,localizations,recordingDetails,snippet,statistics,status,topicDetails"

func (c *Collector) ChannelByUsername(username string) error {
	//kek, err := service.VideoCategories.List([]string{"snippet"}).RegionCode("RU") .Do()
	//fmt.Println(len(kek.Items))

	call := c.service.Channels.List([]string{ChannelParts}).ForUsername(username)
	response, err := call.Do()
	if err != nil {
		return err
	}

	if len(response.Items) == 0 {
		return errors.New("No such user")
	}

	channel := response.Items[0]
	if !db.ChannelExists(channel) {
		err = db.AddChannel(channel)
	}
	db.AddChannelStats(channel)
	//c.VideosByChannelID(channel.Id)
	return err

	//channel := response.Items[0]
	//
	//fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', " +
	//	"and it has %d subscribers and %d videos.",
	//	channel.Id,
	//	channel.Snippet.Title,
	//	channel.Statistics.SubscriberCount, channel.Statistics.VideoCount))

	//err = VideosByChannelID(service, channel.Id)
	//if err != nil {
	//	return err
	//}

	//fmt.Println(response.Items[0].TopicDetails.TopicCategories)
	//fmt.Println(response.Items[0].ContentDetails.RelatedPlaylists)
}

func (c *Collector) ChannelById(id string) error {
	call := c.service.Channels.List([]string{ChannelParts}).Id(id)
	response, err := call.Do()
	if err != nil {
		return err
	}

	channel := response.Items[0]
	if !db.ChannelExists(channel) {
		db.AddChannel(channel)
	}
	db.AddChannelStats(channel)

	return nil
}

func (c *Collector) RelatedVideos(videoId string) error {
	result, err := c.service.Search.List([]string{"snippet"}).RelatedToVideoId(videoId).MaxResults(50).Type("video").Do()
	if err != nil {
		return err
	}
	fmt.Println(result.Items)
	return nil
}

func (c *Collector) VideosByChannelID(channelID string) error {
	search := c.service.Search.List([]string{"id,snippet"}).ChannelId(channelID).PublishedAfter(time.Now().AddDate(0, -1, 0).Format(time.RFC3339)).Order("date").MaxResults(50)
	result, err := search.Do()
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("For the last month: %d videos.", len(result.Items)))

	for _, v := range result.Items {
		vidR, err := c.service.Videos.List([]string{VideoParts}).Id(v.Id.VideoId).Do()
		if err != nil {
			return  err
		}
		//fmt.Println(vidR.Items[0].Statistics.ViewCount)
		if !db.VideoExists(vidR.Items[0]) {
			db.AddVideo(vidR.Items[0])
		}
		db.AddReactions(vidR.Items[0])
	}


	return nil
}
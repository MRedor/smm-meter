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

func (c *Collector) Try() error {
	resp, err := c.service.Search.List([]string{"id,snippet"}).ChannelId("UCDF_NIAEkcAUvzxe1DUzaQA").PublishedAfter(time.Now().AddDate(0, -1, 0).Format(time.RFC3339)).Order("date").MaxResults(50).Do()
	fmt.Println(len(resp.Items))
	return err
}

func (c *Collector) Categories() {
	kek, _ := c.service.VideoCategories.List([]string{"snippet"}).RegionCode("RU") .Do()
	fmt.Println(kek.Items)
}

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

func (c *Collector) ChannelById(id string) (*youtube.Channel, error) {
	call := c.service.Channels.List([]string{ChannelParts}).Id(id)
	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, errors.New("no such channel")
	}

		channel := resp.Items[0]
	if !db.ChannelExists(channel) {
		db.AddChannel(channel)
	}
	db.AddChannelStats(channel)

	return channel, nil
}

func (c *Collector) RelatedVideos(videoId string) error {
	result, err := c.service.Search.List([]string{"snippet"}).RelatedToVideoId(videoId).MaxResults(50).Type("video").Do()
	if err != nil {
		return err
	}
	fmt.Println(result.Items)
	return nil
}

func (c *Collector) VideosForLastMonth(channelId string) ([]*youtube.Video, error) {
	query := c.service.Search.List([]string{"id,snippet"}).ChannelId(channelId).PublishedAfter(time.Now().AddDate(0, -1, 0).Format(time.RFC3339)).Order("date").MaxResults(50)
	resp, err := query.Do()
	if err != nil {
		return nil, err
	}

	ids := ""
	for i, item := range resp.Items {
		if i == 0 {
			ids = item.Id.VideoId
		} else {
			ids = ids + "," + item.Id.VideoId
		}
	}
	videos, err := c.service.Videos.List([]string{VideoParts}).Id(ids).Do()
	if err != nil {
		return nil, err
	}
	return videos.Items, nil
}

func (c *Collector) VideosByChannelID(channelID string) error {
	result, err := c.service.Search.List([]string{"id,snippet"}).ChannelId(channelID).Order("date").MaxResults(50).Do()
	if err != nil {
		return err
	}

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

func (c *Collector) VideoById(id string) (*youtube.Video, error) {
	resp, err := c.service.Videos.List([]string{VideoParts}).Id(id).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Items) == 0 {
		return nil, errors.New("no such video")
	}

	return resp.Items[0], nil
}
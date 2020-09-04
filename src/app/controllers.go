package app

import (
	"collector"
	"db"
	"encoding/json"
	"github.com/labstack/echo"
	"net/http"
)

func GetId(ctx echo.Context) error {
	name := ctx.Param("name")

	c, err :=  collector.NewCollector()
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, err.Error())
	}
	id := c.ChannelByUsername(name)
	return ctx.JSON(http.StatusOK, id)
}

func GetChannel(ctx echo.Context) error {
	id := ctx.Param("id")

	c, err := collector.NewCollector()
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, err.Error())
	}
	channel, err := c.ChannelById(id)
	if err != nil {
		if err.Error() == "no such channel" {
			return ctx.JSON(http.StatusBadRequest, err.Error())
		}
		return ctx.JSON(http.StatusServiceUnavailable, err.Error())
	}

	chResponse := ChannelResponse{
		Title:       channel.Snippet.Title,
		Subscribers: int(channel.Statistics.SubscriberCount),
		VideoCount:  int(channel.Statistics.VideoCount),
	}

	videos, err := c.VideosForLastMonth(id)
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, err.Error())
	}
	chResponse.VideoLastMonth = len(videos)

	var (
		views    = 0
		likes    = 0
		dislikes = 0
		comments = 0
	)
	for _, v := range videos {
		views += int(v.Statistics.ViewCount)
		likes += int(v.Statistics.LikeCount)
		dislikes += int(v.Statistics.DislikeCount)
		comments += int(v.Statistics.CommentCount)
	}
	chResponse.AverageViews = views / len(videos)
	chResponse.AverageLikes = likes / len(videos)
	chResponse.AverageDislikes = dislikes / len(videos)
	chResponse.AverageComments = comments / len(videos)

	if comments != 0 {
		chResponse.SuspiciousVideos = []string{}
		avg := float64(chResponse.AverageComments) / float64(chResponse.AverageViews)

		for _, v := range videos {
			if float64(v.Statistics.CommentCount) / float64(v.Statistics.ViewCount) < avg / 2 {
				chResponse.SuspiciousVideos = append(chResponse.SuspiciousVideos, v.Id)
			}
			if float64(v.Statistics.CommentCount) / float64(v.Statistics.ViewCount) > avg * 5 {
				chResponse.SuspiciousVideos = append(chResponse.SuspiciousVideos, v.Id)
			}
		}
	}


	return ctx.JSON(http.StatusOK, chResponse)
}

func GetVideo(ctx echo.Context) error {
	id := ctx.Param("id")

	c, err := collector.NewCollector()
	if err != nil {
		return err
	}
	video, err := c.VideoById(id)
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, "")
	}

	vidResponse := VideoResponse{
		Title:      video.Snippet.Title,
		Category:   video.Snippet.CategoryId,
		Views:      int(video.Statistics.ViewCount),
		Likes:      int(video.Statistics.LikeCount),
		Dislikes:   int(video.Statistics.DislikeCount),
		Comments:   int(video.Statistics.CommentCount),
		WasPopular: db.VideoWasPopular(video),
	}
	ch, err := c.ChannelById(video.Snippet.ChannelId)
	if err != nil {
		vidResponse.ViewsFromSubscribers = -1
	} else {
		vidResponse.ViewsFromSubscribers = float64(video.Statistics.ViewCount) / float64(ch.Statistics.SubscriberCount)
	}

	return ctx.JSON(http.StatusOK, vidResponse)
}

func GetTimeline(ctx echo.Context) error {
	id := ctx.Param("id")

	stats, err := db.VideoStatsById(id)
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, "")
	}

	timeline := TimeLineResponse{}
	for _, st := range stats {
		timeline.Views = append(timeline.Views, st.Views)
		timeline.Likes = append(timeline.Likes, st.Likes)
		timeline.Dislikes = append(timeline.Dislikes, st.Dislikes)
		timeline.Comments = append(timeline.Comments, st.Comments)
		timeline.Data = append(timeline.Data, st.Data)
	}

	return ctx.JSON(http.StatusOK, timeline)
}

func StartWatching(ctx echo.Context) error {
	req := WatchRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	c, err := collector.NewCollector()
	if err != nil {
		return ctx.JSON(http.StatusServiceUnavailable, "")
	}
	c.WatchVideo(req.Id, req.Period)

	return ctx.JSON(http.StatusOK, "")
}

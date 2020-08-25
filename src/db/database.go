package db

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"

	"config"
)

var db *sqlx.DB

func updateApostrophes(got string) string {
	return strings.Replace(got, "'", "''", -1)
}

func AddVideo(video *youtube.Video) error {
	var contentDetails, recordingDetails, status, topicDetails []byte
	var err error

	if video.ContentDetails != nil {
		contentDetails, err = video.ContentDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	localizations, err := json.Marshal(video.Localizations)
	if err != nil {
		return err
	}

	if video.RecordingDetails != nil {
		recordingDetails, err = video.RecordingDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	snippet, err := video.Snippet.MarshalJSON()
	if err != nil {
		return err
	}

	if video.Status != nil {
		status, err = video.Status.MarshalJSON()
		if err != nil {
			return err
		}
	}

	if video.TopicDetails != nil {
		topicDetails, err = video.TopicDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	query :=  "insert into videos " + "(id, contentDetails, duration, localizations, recordingDetails, snippet, channelId, publishedAt, status, topicDetails)" +
		fmt.Sprintf(` values ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`,
		video.Id, updateApostrophes(string(contentDetails)),
		video.ContentDetails.Duration,
		updateApostrophes(string(localizations)), updateApostrophes(string(recordingDetails)), updateApostrophes(string(snippet)),
		video.Snippet.ChannelId, video.Snippet.PublishedAt,
		updateApostrophes(string(status)), updateApostrophes(string(topicDetails)),
	)
	fmt.Println(query)
	
	_, err = db.Exec(query)

	if err != nil {
		fmt.Println("VIDEOERROR: " + err.Error())
		fmt.Println("query: " + query)
	}

	return err
}

func AddReactions(video *youtube.Video) error {
	query :=  "insert into videoStats (videoId, comments, dislikes, likes, views, favorites)" +
		fmt.Sprintf(` values ('%s', %v, %v, %v, %v, %v)`,
			video.Id, video.Statistics.CommentCount, video.Statistics.DislikeCount, video.Statistics.LikeCount,
			video.Statistics.ViewCount, video.Statistics.FavoriteCount)
	fmt.Println(query)

	_, err := db.Exec(query)
	return err
}

func AddPopular(video *youtube.Video) error {
	query :=  "insert into popular (videoId)" +
		fmt.Sprintf(` values ('%s')`, video.Id)
	fmt.Println(query)

	_, err := db.Exec(query)
	return err
}

func AddChannel(channel *youtube.Channel) error {
	var brandingSettings, contentDetails, contentOwnerDetails, invideoPromotion, status, topicDetails []byte
	var err error

	if channel.BrandingSettings != nil {
		brandingSettings, err = channel.BrandingSettings.MarshalJSON()
		if err != nil {
			return err
		}
	}

	if channel.ContentDetails != nil {
		contentDetails, err = channel.ContentDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	if channel.ContentOwnerDetails != nil {
		contentOwnerDetails, err = channel.ContentOwnerDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	if channel.InvideoPromotion != nil {
		invideoPromotion, err = channel.InvideoPromotion.MarshalJSON()
		if err != nil {
			return err
		}
	}

	localizations, err := json.Marshal(channel.Localizations)
	if err != nil {
		return err
	}

	snippet, err := channel.Snippet.MarshalJSON()
	if err != nil {
		return err
	}

	if channel.Status != nil {
		status, err = channel.Status.MarshalJSON()
		if err != nil {
			return err
		}
	}

	if channel.TopicDetails != nil {
		topicDetails, err = channel.TopicDetails.MarshalJSON()
		if err != nil {
			return err
		}
	}

	query :=  "insert into channels " + "(id, brandingSettings, contentDetails, contentOwnerDetails, invideoPromotion, localizations, snippet, status, topicDetails)" +
		fmt.Sprintf(` values ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')`,
			updateApostrophes(channel.Id), updateApostrophes(string(brandingSettings)),  updateApostrophes(string(contentDetails)),
			updateApostrophes(string(contentOwnerDetails)), updateApostrophes(string(invideoPromotion)),
			updateApostrophes(string(localizations)), updateApostrophes(string(snippet)),
			updateApostrophes(string(status)), updateApostrophes(string(topicDetails)),
		)
	//fmt.Println(query)

	_, err = db.Exec(query)

	if err != nil {
		fmt.Println("CHANNELERROR: " + err.Error())
		fmt.Println("query: " + query)
	}

	return err
}

func AddChannelStats(channel *youtube.Channel) error {
	query :=  "insert into channelStats (channelId, subscriberCount, videoCount)" +
		fmt.Sprintf(` values ('%s', %v, %v)`, channel.Id, channel.Statistics.SubscriberCount, channel.Statistics.VideoCount)
	fmt.Println(query)

	_, err := db.Exec(query)
	return err
}

func VideoExists(video *youtube.Video) bool {
	query := "select count(*) from videos where id='" + video.Id + "'"

	var count int
	row := db.QueryRow(query)
	row.Scan(&count)

	return count > 0
}

func ChannelExists(channel *youtube.Channel) bool {
	query := "select count(*) from channels where id='" + updateApostrophes(channel.Id) + "'"

	var count int
	row := db.QueryRow(query)
	row.Scan(&count)

	return count > 0
}

func VideoWasPopular(video *youtube.Video) bool {
	query := "select count(*) from popular where videoId='" + video.Id + "'"

	var count int
	row := db.QueryRow(query)
	row.Scan(&count)

	return count > 0
}

type VideoStats struct {
	Id string `db:"videoId"`
	Comments int `db:"comments"`
	Dislikes int `db:"dislikes"`
	Likes int `db:"likes"`
	Views int `db:"views"`
	Favorites int `db:"favorites"`
	Data string `db:"data"`
}

func VideoStatsById(videoId string) ([]VideoStats, error) {
	query := fmt.Sprintf("select * from videoStats where convert(varchar, videoId)='%s' order by data", videoId)

	stats := []VideoStats{}
	err := db.Select(&stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func init() {
	source := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		config.Cfg.DB.Host, config.Cfg.DB.Username, config.Cfg.DB.Password, config.Cfg.DB.Port, config.Cfg.DB.Database)
	database, err := sqlx.Open("sqlserver", source)
	if err != nil {
		log.Fatal(err)
	}
	db = database
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	//query := "select count(*) from videos where id='" + "kek" + "'"
	//
	//var count int
	//row := db.QueryRow(query)
	//row.Scan(&count)
	//
	//fmt.Println(count)
}
package app

type ChannelResponse struct {
	Title            string   `json:"title"`
	Subscribers      int      `json:"subscribers"`
	VideoCount       int      `json:"videoCount"`
	VideoLastMonth   int      `json:"videoLastMonth"`
	AverageViews     int      `json:"averageViews"`
	AverageLikes     int      `json:"averageLikes"`
	AverageDislikes  int      `json:"averageDislikes"`
	AverageComments  int      `json:"averageComments"`
	SuspiciousVideos []string `json:"suspiciousVideos"`
}

type VideoResponse struct {
	Title                string  `json:"title"`
	Category             string  `json:"category"`
	Views                int     `json:"views"`
	Likes                int     `json:"likes"`
	Dislikes             int     `json:"dislikes"`
	Comments             int     `json:"comments"`
	WasPopular           bool    `json:"wasPopular"`
	ViewsFromSubscribers float64 `json:"viewsFromSubscribers"`
}

type TimeLineResponse struct {
	Views    []int    `json:"views"`
	Likes    []int    `json:"likes"`
	Dislikes []int    `json:"dislikes"`
	Comments []int    `json:"comments"`
	Data     []string `json:"data"`
}

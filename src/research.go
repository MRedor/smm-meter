package main

import (
	"bufio"
	"collector"
	"db"
	"golang.org/x/exp/errors/fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func plotPoints(p []float64) plotter.XYs {
	pts := make(plotter.XYs, len(p))
	for i := range pts {
		pts[i].Y = p[i]
		pts[i].X = 0
	}
	return pts
}

func CommentsViews()  {
	file1, err := os.Open("src/data/comments")
	if err != nil {
		fmt.Println(err)
	}
	defer file1.Close()
	reader1 := bufio.NewReader(file1)

	file2, err := os.Open("src/data/views")
	if err != nil {
		fmt.Println(err)
	}
	defer file2.Close()
	reader2 := bufio.NewReader(file2)

	var res []float64
	var smalles []float64
	for {
		line1, _, err := reader1.ReadLine()
		if err != nil {
			break
		}
		line2, _, err := reader2.ReadLine()
		if err != nil {
			break
		}

		int1, _ := strconv.Atoi(string(line1))
		int2, _ := strconv.Atoi(string(line2))
		val := float64(int1)/float64(int2)
		res = append(res, val)
		if val < 0.01 {
			smalles = append(smalles, val)
		}

	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(p, "", plotPoints(res))
	if err != nil {
		panic(err)
	}
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "src/data/commentsviews.png"); err != nil {
		panic(err)
	}

	pSmall, err := plot.New()
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(pSmall, "", plotPoints(smalles))
	if err != nil {
		panic(err)
	}
	if err := pSmall.Save(5*vg.Inch, 5*vg.Inch, "src/data/commentsviewsSmall.png"); err != nil {
		panic(err)
	}
}

func ReactionsViews() {
	fileLikes, err := os.Open("src/data/likes")
	if err != nil {
		fmt.Println(err)
	}
	defer fileLikes.Close()
	readerLikes := bufio.NewReader(fileLikes)

	fileDislikes, err := os.Open("src/data/dislikes")
	if err != nil {
		fmt.Println(err)
	}
	defer fileDislikes.Close()
	readerDislikes := bufio.NewReader(fileDislikes)

	fileViews, err := os.Open("src/data/views")
	if err != nil {
		fmt.Println(err)
	}
	defer fileViews.Close()
	readerViews := bufio.NewReader(fileViews)

	var res []float64
	for {
		lineLikes, _, err := readerLikes.ReadLine()
		if err != nil {
			break
		}
		lineDislikes, _, err := readerDislikes.ReadLine()
		if err != nil {
			break
		}
		lineViews, _, err := readerViews.ReadLine()
		if err != nil {
			break
		}

		intLikes, _ := strconv.Atoi(string(lineLikes))
		intDislikes, _ := strconv.Atoi(string(lineDislikes))
		intViews, _ := strconv.Atoi(string(lineViews))
		val := float64(intLikes + intDislikes)/float64(intViews)
		if val < 0.02 {
			res = append(res, val)
		}
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(p, "", plotPoints(res))
	if err != nil {
		panic(err)
	}
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "src/data/reactionsviews.png"); err != nil {
		panic(err)
	}
}

func DislikesLikes() {
	fileLikes, err := os.Open("src/data/likes")
	if err != nil {
		fmt.Println(err)
	}
	defer fileLikes.Close()
	readerLikes := bufio.NewReader(fileLikes)

	fileDislikes, err := os.Open("src/data/dislikes")
	if err != nil {
		fmt.Println(err)
	}
	defer fileDislikes.Close()
	readerDislikes := bufio.NewReader(fileDislikes)

	var res []float64
	for {
		lineLikes, _, err := readerLikes.ReadLine()
		if err != nil {
			break
		}
		lineDislikes, _, err := readerDislikes.ReadLine()
		if err != nil {
			break
		}

		intLikes, _ := strconv.Atoi(string(lineLikes))
		intDislikes, _ := strconv.Atoi(string(lineDislikes))
		if intLikes == 0 {
			continue
		}
		val := float64(intDislikes)/float64(intLikes)
		res = append(res, val)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(p, "", plotPoints(res))
	if err != nil {
		panic(err)
	}
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "src/data/dislikesLikes.png"); err != nil {
		panic(err)
	}
}

func SubscriberViews() {
	file, err := os.Open("src/data/subscriberviews")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	var res []float64
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		kek := strings.Split(string(line), ",")

		intSub, _ := strconv.Atoi(kek[0])
		intViews, _ := strconv.Atoi(kek[1])
		val := float64(intViews)/float64(intSub)
		if val > 0.25 {
			continue
		}
		res = append(res, val)

	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}
	err = plotutil.AddLinePoints(p, "", plotPoints(res))
	if err != nil {
		panic(err)
	}
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "src/data/subscriberviews.png"); err != nil {
		panic(err)
	}
}

func TimeFromData(data string) float64 {
	h, _ := strconv.Atoi(data[11:13])
	m, _ := strconv.Atoi(data[14:16])
	return float64(h) + float64(m) / 60
}

func commentsPoints(stats []db.VideoStats) plotter.XYs {
	pts := make(plotter.XYs, len(stats))
	for i := range pts {
		pts[i].Y = float64(stats[i].Comments)
		pts[i].X = TimeFromData(stats[i].Data)
	}
	return pts
}

func likesPoints(stats []db.VideoStats) plotter.XYs {
	pts := make(plotter.XYs, len(stats))
	for i := range pts {
		if i == 0 {
			pts[i].Y = 0
			pts[i].X = TimeFromData(stats[i].Data)
			continue
		}
		pts[i].Y = float64(stats[i].Likes - stats[i - 1].Likes) / 1000
		pts[i].X = TimeFromData(stats[i].Data)
	}
	return pts
}

func dislikesPoints(stats []db.VideoStats) plotter.XYs {
	pts := make(plotter.XYs, len(stats))
	for i := range pts {
		if i == 0 {
			pts[i].Y = 0
			pts[i].X = TimeFromData(stats[i].Data)
			continue
		}
		pts[i].Y = float64(stats[i].Dislikes - stats[i - 1].Dislikes)
		pts[i].X = TimeFromData(stats[i].Data)
	}
	return pts
}

func reactionsPoints(stats []db.VideoStats) plotter.XYs {
	pts := make(plotter.XYs, len(stats))
	for i := range pts {
		if i == 0 {
			pts[i].Y = 0
			pts[i].X = TimeFromData(stats[i].Data)
			continue
		}
		pts[i].Y = float64(stats[i].Likes + stats[i].Dislikes - stats[i - 1].Likes - stats[i - 1].Dislikes) / 1000
		pts[i].X = TimeFromData(stats[i].Data)
	}
	return pts
}

func viewsPoints(stats []db.VideoStats) plotter.XYs {
	pts := make(plotter.XYs, len(stats))
	for i := range pts {
		if i == 0 {
			pts[i].Y = 0
			pts[i].X = TimeFromData(stats[i].Data)
			continue
		}
		pts[i].Y = float64(stats[i].Views - stats[i - 1].Views) / 10000
		pts[i].X = TimeFromData(stats[i].Data)
	}
	return pts
}

func VideoPlot(videoId string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Video:" + videoId + " timeline"
	p.X.Label.Text = "time"
	p.Y.Label.Text = "number"

	stats, _ := db.VideoStatsById(videoId)
	if len(stats) < 5 {

	}

		err = plotutil.AddLinePoints(p,
		//"comments", commentsPoints(stats),
		"dislikes+likes", reactionsPoints(stats),
		"likes", likesPoints(stats),
		//"dislikes", dislikesPoints(stats),
		"views", viewsPoints(stats))

	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(5*vg.Inch, 5*vg.Inch, "src/plots/" + videoId + ".png"); err != nil {
		panic(err)
	}
}

func TrendingVideosPlots() {
	file, err := os.Open("src/collector/trends_on_watch")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		VideoPlot(string(line))
	}
}

func AvgtrendingTime() {
	file, err := os.Open("src/collector/trends_on_watch")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	count := float64(0)
	sum := float64(0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		stats, _ := db.VideoStatsById(string(line))
		if len(stats) < 5 {
			continue
		}
		count = count + 1
		sum = sum + TimeFromData(stats[len(stats) - 1].Data) - TimeFromData(stats[0].Data)
	}
	fmt.Println(sum / count)
}

func th1() {
	f, err := os.OpenFile("src/category_th", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()


	file, err := os.Open("src/collector/channels_on_watch")
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
		//string(line)
		videos, err := c.VideosForLastMonth(string(line))
		if err != nil {
			fmt.Println(err.Error())
		}
		if len(videos) > 1 {
			fmt.Println("HEY")
		}
		//slice.SortInterface(videos[:], func(i, j int) bool {
		//	return videos[i].Statistics.ViewCount < videos[j].Statistics.ViewCount
		//})

		for _, v := range videos {
			//fmt.Println(, " ", )
			if _, err := f.WriteString(strconv.Itoa(int(v.Statistics.ViewCount)) + " " + v.Snippet.CategoryId + "\n"); err != nil {
				log.Println(err)
			}
		}

		if _, err := f.WriteString("---------------------------------\n"); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	//db.VideoStatsById("gdZLi9oWNZg")
	//TrendingVideosPlots()

	//AvgtrendingTime()

	//c, _ := collector.NewCollector()
	//c.Categories()

	//th1()
	//CommentsViews()
	//ReactionsViews()
	//DislikesLikes()
	SubscriberViews()
}

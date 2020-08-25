package main

import (
	"bufio"
	"collector"
	"db"
	"golang.org/x/exp/errors/fmt"
	"os"
	"strconv"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

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

func main() {
	//db.VideoStatsById("gdZLi9oWNZg")
	//TrendingVideosPlots()

	//AvgtrendingTime()
	c, _ := collector.NewCollector()
	c.Categories()
}

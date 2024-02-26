package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func TimeSeries(t [][2]any) {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic line example", Subtitle: "This is the subtitle."}),
	)

	ts := make([]any, 0, len(t))
	vs := make([]opts.LineData, 0, len(t))
	for i := 0; i < len(t); i++ {
		ts = append(ts, t[i][0])
		vs = append(vs, opts.LineData{
			Value: t[i][1],
		})
	}

	line.SetXAxis(ts).
		AddSeries("line", vs)

	page := components.NewPage()
	page.AddCharts(
		line,
	)
	f, err := os.Create("web/html/line.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}

package main

import (
	"fmt"
	"image/color"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dpdf"
)

const A4_W float64 = 210.0
const A4_H float64 = 297.0
const MARGIN_W float64 = 10.0
const MARGIN_H float64 = 18.0 // I don't know why, but keeping a smaller bottom margin page breaks "365" on last line...
const EXTRA_TOP_MARGIN float64 = 7.0

const COLS = 14
const ROWS = 27 // ceil(365 / 14)

const START_X = MARGIN_W
const END_X = A4_W - MARGIN_W
const START_Y = MARGIN_H + EXTRA_TOP_MARGIN
const END_Y = A4_H - MARGIN_H

const BOX_W = (END_X - START_X) / float64(COLS)
const BOX_H = (END_Y - START_Y) / float64(ROWS)
const BOX_HALF_H = BOX_H / 2
const BOX_HALF_W = BOX_W / 2

func generatePdf() (*gofpdf.Fpdf, error) {
	draw2d.SetFontFolder("./resource/font/")

	// Initialize the graphic context on an RGBA image
	dest := draw2dpdf.NewPdf("P", "mm", "A4")
	gc := draw2dpdf.NewGraphicContext(dest)

	calStartTimeStr := "2021-08-02T00:00:00Z" // set this to whatever you want...
	calStartTime, err := time.Parse(time.RFC3339, calStartTimeStr)
	if err != nil {
		return nil, err
	}

	//interFontMedium := draw2d.FontData{Name: "Inter-Medium", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}
	interFontSemiBold := draw2d.FontData{Name: "Inter-SemiBold", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}
	interFontBold := draw2d.FontData{Name: "Inter-Bold", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}
	interFontExtraBold := draw2d.FontData{Name: "Inter-ExtraBold", Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}

	// Draw heading "Seinfeld calendar" + year/start date
	gc.Save()
	heading := "Seinfeld calendar"
	subheading1 := calStartTime.Format("2006")
	subheading2 := fmt.Sprintf("Start: %s", calStartTime.Format("Mon January 2"))
	gc.SetFontData(interFontExtraBold)
	gc.SetFontSize(10)
	gc.SetFillColor(color.RGBA{0xff, 0x0, 0x0, 0xff})
	gc.FillStringAt(heading, START_X, START_Y-6)
	gc.SetFontData(interFontBold)
	gc.SetFontSize(3)
	gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	gc.FillStringAt(subheading1, 106, START_Y-10)
	gc.FillStringAt(subheading2, 106, START_Y-6)
	gc.Restore()

	gc.Save()
	gc.SetFillColor(color.Transparent)
	gc.SetStrokeColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	gc.SetLineWidth(0.25)

	// Draw col lines:
	for i := 0; i < COLS+1; i++ {
		x := (END_X-START_X)/float64(COLS)*float64(i) + START_X
		gc.MoveTo(x, START_Y)
		gc.LineTo(x, END_Y)
		gc.Stroke()
	}

	// Draw row lines:
	for i := 0; i < ROWS+1; i++ {
		y := (END_Y-START_Y)/float64(ROWS)*float64(i) + START_Y
		gc.MoveTo(START_X, y)
		gc.LineTo(END_X, y)
		gc.Stroke()
	}

	gc.Restore()
	gc.Save()

	gc.SetFontData(interFontBold)
	gc.SetFontSize(3.5)
	gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	gc.SetStrokeColor(color.Black)

	// Print numbers:
	for i := 0; i < 365; i++ {
		gc.SetFontData(interFontBold)

		date := calStartTime.AddDate(0, 0, i)
		dayNbr := fmt.Sprint(i + 1)

		col := i % COLS
		row := int(math.Ceil(float64(i+1)/float64(COLS))) - 1

		// Set x, y as top left of box:
		boxStartX := (END_X-START_X)/float64(COLS)*float64(col) + START_X
		boxStartY := (END_Y-START_Y)/float64(ROWS)*float64(row) + START_Y

		_, tt, tr, _ := gc.GetStringBounds(dayNbr)

		// Put text centered on half of box width/height
		x := boxStartX + BOX_HALF_W
		y := boxStartY + BOX_HALF_H
		// With some small adjustments for text width/height
		y += math.Abs(tt) / 2
		x -= math.Abs(tr) / 2
		// Adjustment for subtext:
		y -= 0.5

		gc.FillStringAt(dayNbr, x, y)

		// Print date of month:
		gc.Save()
		gc.SetFontData(interFontSemiBold)
		gc.SetFontSize(1.5)
		gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
		subtext := date.Format("2")
		if date.Day() == 1 {
			subtext = date.Format("Jan")
		}
		subtext += " " + date.Weekday().String()[0:1]
		// Make start of month & weekends stand out:
		// if date.Day() == 1 || date.Weekday() == 0 || date.Weekday() == 6 {
		// 	gc.SetFillColor(color.Black)
		// }
		_, _, etr, _ := gc.GetStringBounds(subtext)
		// this is an ugly solution for some weird page breaking on day "365"
		// could also be solved with an offset but then it will overlap...
		if dayNbr != "365" {
			gc.FillStringAt(subtext, boxStartX+(BOX_W/2)-(etr/2), boxStartY+BOX_H-1.5)
		}
		gc.Restore()
	}

	gc.Restore()

	return dest, nil
}

func handlePdf(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Handling request...")

	result, err := generatePdf()

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Something went wrong.", 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=\"calendar.pdf\"")

	err = result.Output(w)

	if err != nil {
		fmt.Println("Error:", err)
		http.Error(w, "Something went wrong.", 500)
		return
	}
}

func main() {
	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	http.HandleFunc("/pdf", handlePdf)
	http.ListenAndServe(":"+port, nil)
}

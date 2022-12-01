package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"

	"gonum.org/v1/gonum/interp"
)

func main() {
	fmt.Println("Rendering colormap...")

	cmap_name := "mycolormap"

	n_color := 1024

	xs := []float64{0.0, 0.16, 0.42, 0.6425, 0.8575, 1.0}
	yrs := []float64{0, 32, 237, 255, 0, 0}
	ygs := []float64{7, 107, 255, 170, 2, 7}
	ybs := []float64{100, 203, 255, 0, 0, 100}

	cmap_data := make([]color.NRGBA, 0, n_color)

	var plr interp.PiecewiseLinear
	var plg interp.PiecewiseLinear
	var plb interp.PiecewiseLinear

	plr.Fit(xs, yrs)
	plg.Fit(xs, ygs)
	plb.Fit(xs, ybs)
	for i := 0; i < n_color; i++ {
		scaled_i := float64(i) / float64(n_color)
		r := plr.Predict(scaled_i)
		g := plg.Predict(scaled_i)
		b := plb.Predict(scaled_i)
		a := 255
		cmap_data = append(cmap_data, color.NRGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: uint8(a),
		})
	}
	save_cmap_to_csv(cmap_data, cmap_name)
	save_cmap_to_png(cmap_data, cmap_name)

	fmt.Println("Data saved to file", cmap_name)
}

func save_cmap_to_csv(cmap_data []color.NRGBA, cmap_name string) {
	color_map_file, _ := os.Create(fmt.Sprintf("./%v.csv", cmap_name))
	defer color_map_file.Close()
	w := csv.NewWriter(color_map_file)
	defer w.Flush()
	header := make([]string, 0, 4)
	header = append(header, "R", "G", "B", "A")
	w.Write(header)
	for _, data := range cmap_data {
		line := []string{strconv.Itoa(int(data.R)), strconv.Itoa(int(data.G)), strconv.Itoa(int(data.B)), strconv.Itoa(int(data.A))}
		if err := w.Write(line); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
}

func save_cmap_to_png(cmap_data []color.NRGBA, cmap_name string) {
	height := 100
	width := len(cmap_data)
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for i, data := range cmap_data {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.NRGBA{
				R: uint8(data.R),
				G: uint8(data.G),
				B: uint8(data.B),
				A: uint8(data.A),
			})
		}
	}
	f, _ := os.Create(fmt.Sprintf("./%v.png", cmap_name))
	defer f.Close()
	png.Encode(f, img)
	f.Close()
}

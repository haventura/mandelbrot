package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/gonum/interp"
)

func main() {
	// go run . -path="..\output_data\output_2022_11_18-23_44_13.csv"
	fmt.Printf("Rendering...\n")
	start_rend := time.Now()

	var path_flag = flag.String("path", "../output_data/output.csv", "Relative path of the csv file.")
	flag.Parse()
	dataset, err := readData(*path_flag)
	s := strings.Split(*path_flag, "\\")
	img_name := strings.Split(s[len(s)-1], ".")[0]
	width, _ := strconv.Atoi(dataset[len(dataset)-1][0])
	width += 1
	height, _ := strconv.Atoi(dataset[len(dataset)-1][1])
	height += 1
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	colormap := compute_colormap(1024)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range dataset {

		pixel_x, _ := strconv.Atoi(line[0])
		pixel_y, _ := strconv.Atoi(line[1])
		it_escape, _ := strconv.Atoi(line[2])
		r_part, _ := strconv.ParseFloat(line[3], 64)
		im_part, _ := strconv.ParseFloat(line[4], 64)
		mod, _ := strconv.ParseFloat(line[5], 64)
		if mod < 2 {
			//pts = append(pts, plotter.XY{X: x, Y: y})
			img.Set(pixel_x, pixel_y, color.NRGBA{
				R: uint8(0),
				G: uint8(0),
				B: uint8(0),
				A: 255,
			})
		} else {
			img.Set(pixel_x, pixel_y, get_color(it_escape, complex(r_part, im_part), colormap))
		}
	}

	f, err := os.Create(fmt.Sprintf("../output_images/%v.png", img_name))
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	elapsed_rend := time.Since(start_rend)
	fmt.Printf("Rendering took %s", elapsed_rend)
}

func readData(fileName string) ([][]string, error) {

	f, err := os.Open(fileName)

	if err != nil {
		return [][]string{}, err
	}

	defer f.Close()

	r := csv.NewReader(f)

	// skip first line
	if _, err := r.Read(); err != nil {
		return [][]string{}, err
	}

	records, err := r.ReadAll()

	if err != nil {
		return [][]string{}, err
	}

	return records, nil
}

func compute_colormap(n_color int) []color.NRGBA {

	h := 100
	img := image.NewNRGBA(image.Rect(0, 0, n_color, h))

	xs := []float64{0.0, 0.16, 0.42, 0.6425, 0.8575, 1.0}
	yrs := []float64{0, 32, 237, 255, 0, 0}
	ygs := []float64{7, 107, 255, 170, 2, 7}
	ybs := []float64{100, 203, 255, 0, 0, 100}

	colormap := make([]color.NRGBA, 0, n_color)

	var plr interp.PiecewiseLinear
	var plg interp.PiecewiseLinear
	var plb interp.PiecewiseLinear
	// 	0.0: color.NRGBA{R: uint8(0), G: uint8(7), B: uint8(100), A: 255,},
	// 	0.16: color.NRGBA{R: uint8(32), G: uint8(107), B: uint8(203), A: 255,},
	// 	0.42: color.NRGBA{R: uint8(237), G: uint8(255), B: uint8(255), A: 255,},
	// 	0.6425: color.NRGBA{R: uint8(255), G: uint8(170), B: uint8(0), A: 255,},
	// 	0.8575: color.NRGBA{R: uint8(0), G: uint8(2), B: uint8(0), A: 255,},
	plr.Fit(xs, yrs)
	plg.Fit(xs, ygs)
	plb.Fit(xs, ybs)
	for i := 0; i < n_color; i++ {
		scaled_i := float64(i) / float64(n_color)
		r := plr.Predict(scaled_i)
		g := plg.Predict(scaled_i)
		b := plb.Predict(scaled_i)
		colormap = append(colormap, color.NRGBA{
			R: uint8(r),
			G: uint8(g),
			B: uint8(b),
			A: 255,
		})
		for j := 0; j < h; j++ {
			img.Set(i, j, color.NRGBA{
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: 255,
			})
		}
	}
	f, err := os.Create("colormap.png")
	if err != nil {
		log.Fatal(err)
	}
	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	return colormap
}

func get_color(it int, c complex128, colors []color.NRGBA) color.NRGBA {
	//found on https://stackoverflow.com/questions/16500656/which-color-gradient-is-used-to-color-mandelbrot-in-wikipedia
	smoothed := math.Log2(math.Log2(math.Pow(real(c), 2)+math.Pow(imag(c), 2)) / 2)
	colorI := (int)(math.Sqrt(float64(it)+10.0-smoothed)*256) % len(colors)

	color := colors[colorI]
	return color
}

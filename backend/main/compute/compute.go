package compute

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"mandelbrot/andre/main/model"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

type point struct {
	x         int
	y         int
	complex   complex128
	escape_it int
}

func Compute(data model.Compute_data) {
	current_time := time.Now()
	file_name := fmt.Sprintf("output_%d_%02d_%02d-%02d_%02d_%02d",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second())

	width := data.Width
	height := data.Height

	//const width, height = 1024, 1024
	//const max_iteration = 128
	n := width * height
	data_c := make(chan point, n)

	// // Reverse seahorse
	// const center_r = -0.743030
	// const center_i = 0.126433
	// const radius = 0.009
	// const min_r, max_r = center_r - radius, center_r + radius
	// const min_i, max_i = center_i - radius, center_i + radius
	// //img full
	// const min_r, max_r = -2.0, 0.5
	// const min_i, max_i = -1.25, 1.25
	// //img 1
	// const min_r, max_r = -2.0, -1.5
	// const min_i, max_i = -0.25, 0.25
	//"mycolormap"

	fmt.Printf("Computing...\n")
	wg.Add(n)
	start_comp := time.Now()
	for x := 0; x < width; x++ {
		x := x
		for y := 0; y < height; y++ {
			y := y
			go worker(x, y, width, height, data.Min_r, data.Max_r, data.Min_i, data.Max_i, data.Max_iteration, data_c)
		}
	}
	wg.Wait()
	close(data_c)
	elapsed_comp := time.Since(start_comp)
	fmt.Printf("Computing %v points over %v iterations took %s\n", n, data.Max_iteration, elapsed_comp)

	//save_to_csv(file_name, data_c)

	start_rend := time.Now()
	render(width, height, data_c, file_name, data.Colormap_name)
	elapsed_rend := time.Since(start_rend)
	fmt.Printf("Rendering took %s", elapsed_rend)
}

func worker(x int, y int, width int, height int, min_r float64, max_r float64, min_i float64, max_i float64, max_iteration int, data_c chan<- point) {
	defer wg.Done()
	r := scale_px_to_coord(x, width-1, min_r, max_r)
	i := scale_px_to_coord(y, height-1, min_i, max_i)
	c := complex(r, i)
	z, it := mandelbrot(c, max_iteration)

	data := point{x, y, z, it}
	data_c <- data
}

func scale_px_to_coord(im_val int, im_max int, mend_min float64, mend_max float64) float64 {
	//        (b-a)(x - min)
	// f(x) = --------------  + a
	// 		     max - min

	scaled := ((mend_max - mend_min) * float64(im_val) / float64(im_max)) + mend_min

	return scaled
}

func mandelbrot(c complex128, max_iteration int) (complex128, int) {
	var z complex128 = complex(0, 0)
	var mod float64
	n := 0
	for i := 0; i < max_iteration; i++ {
		n++
		z = cmplx.Pow(z, 2) + c
		mod = cmplx.Abs(z)
		if mod > 2 {
			break
		}
	}
	return z, n
}

func render(width int, height int, dataset chan point, img_name string, cmap_name string) {

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	colormap := read_cmap_from_csv(fmt.Sprintf("./colormaps/%v.csv", cmap_name))
	for p := range dataset {
		if cmplx.Abs(p.complex) < 2 {
			img.Set(p.x, p.y, color.NRGBA{
				R: uint8(0),
				G: uint8(0),
				B: uint8(0),
				A: 255,
			})
		} else {
			img.Set(p.x, p.y, get_color_from_cmap(p.complex, p.escape_it, colormap))
		}
	}
	f, err := os.Create(fmt.Sprintf("./output_images/%v.png", img_name))
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
}

func get_color_from_cmap(c complex128, it int, cmap []color.NRGBA) color.NRGBA {
	//found on https://stackoverflow.com/questions/16500656/which-color-gradient-is-used-to-color-mandelbrot-in-wikipedia
	smoothed := math.Log2(math.Log2(math.Pow(real(c), 2)+math.Pow(imag(c), 2)) / 2)

	color_index := (int)(math.Sqrt(float64(it)+10.0-smoothed)*256) % len(cmap)
	color := cmap[color_index]
	return color
}

func read_cmap_from_csv(cmap_name string) []color.NRGBA {
	var cmap_data []color.NRGBA
	f, _ := os.Open(cmap_name)
	defer f.Close()
	r := csv.NewReader(f)
	// skip first line
	r.Read()
	records, _ := r.ReadAll()
	for _, record := range records {
		R, _ := strconv.Atoi(record[0])
		G, _ := strconv.Atoi(record[1])
		B, _ := strconv.Atoi(record[2])
		A, _ := strconv.Atoi(record[3])
		cmap_data = append(cmap_data, color.NRGBA{uint8(R), uint8(G), uint8(B), uint8(A)})
	}
	return cmap_data
}

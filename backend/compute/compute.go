package compute

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/png"
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

func Compute_image(data model.Image_data) []byte {
	width := data.Width
	height := data.Height

	n := width * height
	data_c := make(chan point, n)

	fmt.Printf("Computing...\n")
	wg.Add(n)
	start_comp := time.Now()
	for x := 0; x < width; x++ {
		x := x
		for y := 0; y < height; y++ {
			y := y
			go worker(x, y, 0, 0, width-1, height-1, data.Min_r, data.Max_r, data.Min_i, data.Max_i, data.Max_iteration, data_c)
		}
	}
	wg.Wait()
	close(data_c)
	elapsed_comp := time.Since(start_comp)
	fmt.Printf("Computing %v points over %v iterations took %s\n", n, data.Max_iteration, elapsed_comp)

	start_rend := time.Now()
	image := render(width, height, data_c, data.Colormap_name)
	elapsed_rend := time.Since(start_rend)
	fmt.Printf("Rendering took %s", elapsed_rend)
	return image
}

func worker(x int, y int, x_min int, y_min, x_max int, y_max int, r_min float64, r_max float64, i_min float64, i_max float64, max_iteration int, data_c chan<- point) {
	defer wg.Done()
	r := scale_px_to_coord(x, x_min, x_max, r_min, r_max)
	i := scale_px_to_coord(y, y_min, y_max, i_max, i_min) // !!! y-axis direction of pixels and complex-plane are inverted !!!
	c := complex(r, i)
	z, it := mandelbrot(c, max_iteration)

	data := point{x, y, z, it}
	data_c <- data
}

func scale_px_to_coord(pixel_value int, pixel_minimum int, pixel_maximum int, coord_minimum float64, coord_maximum float64) float64 {
	//        (b-a)(x - min)
	// f(x) = --------------  + a
	// 		     max - min

	scaled := ((coord_maximum - coord_minimum) * float64(pixel_value-pixel_minimum) / float64(pixel_maximum-pixel_minimum)) + coord_minimum

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

func render(width int, height int, dataset chan point, cmap_name string) []byte {
	my_image := image.NewNRGBA(image.Rect(0, 0, width, height))
	colormap := read_cmap_from_csv(fmt.Sprintf("./colormap/%v.csv", cmap_name))
	for p := range dataset {
		if cmplx.Abs(p.complex) < 2 {
			my_image.Set(p.x, p.y, color.NRGBA{
				R: uint8(0),
				G: uint8(0),
				B: uint8(0),
				A: 255,
			})
		} else {
			my_image.Set(p.x, p.y, get_color_from_cmap(p.complex, p.escape_it, colormap))
		}
	}
	buffer := new(bytes.Buffer)
	png.Encode(buffer, my_image)
	image_bytes := buffer.Bytes()
	return image_bytes
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

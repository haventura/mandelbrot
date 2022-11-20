package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"
	"sync"
	"time"
)

func main() {

	current_time := time.Now()
	file_name := fmt.Sprintf("output_%d_%02d_%02d-%02d_%02d_%02d",
		current_time.Year(), current_time.Month(), current_time.Day(),
		current_time.Hour(), current_time.Minute(), current_time.Second())
	const width, height = 1024, 1024
	n := width * height

	const max_it = 128

	// // Reverse seahorse
	const center_r = -0.743030
	const center_i = 0.126433
	const diameter = 0.009
	const min_r, max_r = center_r - diameter, center_r + diameter
	const min_i, max_i = center_i - diameter, center_i + diameter
	// //img full
	// const min_r, max_r = -2.0, 0.5
	// const min_i, max_i = -1.25, 1.25
	// //img 1
	// const min_r, max_r = -2.0, -1.5
	// const min_i, max_i = -0.25, 0.25

	dataset := make([][]string, 0, n+1)
	header := make([]string, 0, 6)
	header = append(header, "pixel_x", "pixel_y", "it_escape", "r_part", "im_part")
	data_c := make(chan []string, n+1)
	data_c <- header
	var wg sync.WaitGroup
	wg.Add(n)
	fmt.Printf("Computing...\n")
	start_comp := time.Now()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			x, y := x, y
			go func() {
				worker(x, y, width, height, min_r, max_r, min_i, max_i, max_it, data_c)
				for i := 0; i < len(data_c); i++ {
					dataset = append(dataset, <-data_c)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
	elapsed_comp := time.Since(start_comp)
	fmt.Printf("Computing %v points over %v iterations took %s\n", n, max_it, elapsed_comp)

	data_file, err := os.Create(fmt.Sprintf("../output_data/%v.csv", file_name))
	if err != nil {
		log.Fatalln("failed to open file", err)
	}
	defer data_file.Close()
	w := csv.NewWriter(data_file)
	defer w.Flush()
	for _, line := range dataset {
		if err := w.Write(line); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}
	fmt.Printf("Data saved to file %v\n", file_name)
}

func worker(x int, y int, width int, height int, min_r float64, max_r float64, min_i float64, max_i float64, max_it int, data_c chan<- []string) {
	r := scale_px_to_coord(x, width-1, min_r, max_r)
	i := scale_px_to_coord(y, height-1, min_i, max_i)
	r, i, it := mandelbrot(r, i, max_it)

	data := make([]string, 0, 6)
	data = append(data, fmt.Sprint(x), fmt.Sprint(y), fmt.Sprint(it), fmt.Sprintf("%.6g", r), fmt.Sprintf("%.6g", i))
	data_c <- data
}

func scale_px_to_coord(im_val int, im_max int, mend_min float64, mend_max float64) float64 {
	//        (b-a)(x - min)
	// f(x) = --------------  + a
	// 		     max - min

	scaled := ((mend_max - mend_min) * float64(im_val) / float64(im_max)) + mend_min

	return scaled
}

func mandelbrot(r_part float64, i_part float64, max_iteration int) (float64, float64, int) {
	c := complex(r_part, i_part)
	var z complex128 = complex(0, 0)
	var mod float64
	n := 0

	for i := 0; i < max_iteration; i++ {
		n++
		z = cmplx.Pow(z, 2) + c
		mod = math.Sqrt(math.Pow(real(z), 2) + math.Pow(imag(z), 2))
		if mod > 2 {
			break
		}
	}
	return real(z), imag(z), n
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
	"time"

	"gonum.org/v1/gonum/interp"
)

func main(){
	
	img_name := "output_200_it"
	const width, height = 900, 900
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	colormap := compute_colormap(256)
	const min_r, max_r = -2.0, 0.5
	const min_i, max_i = -1.25, 1.25
	const max_it = 200
	//step_r := (max_r - min_r) / width
	//step_i := (max_i - min_i) / height

	n := width * height

	fmt.Printf("Computing...\n")
	start_comp := time.Now()
	for x := 0; x < width; x ++  {
		for y := 0; y < height; y ++ {
			r := scale_px_to_coord(x, width, min_r, max_r)
			i := scale_px_to_coord(y, height, min_i, max_i)
			val, it := mandelbrot(r,i,max_it)
			if val < 2 {
				//pts = append(pts, plotter.XY{X: x, Y: y})
				img.Set(x, y, color.NRGBA{
					R: uint8(0),
					G: uint8(0),
					B: uint8(0),
					A: 255,
				})
			} else {
				img.Set(x, y, get_color(it, complex(r,i), colormap))
			}
			//fmt.Printf("(%v, %v), val: %v, it: %v\n", x, y, val, it)
		}
	}
	elapsed_comp := time.Since(start_comp)
	fmt.Printf("computing %v points over %v iterations took %s\n", n, max_it, elapsed_comp)
	fmt.Printf("Rendering...\n")
	start_rend := time.Now()
	f, err := os.Create(fmt.Sprintf("%v.png", img_name))
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

func scale_px_to_coord(im_val int, im_max int, mend_min float64, mend_max float64) float64{
	//        (b-a)(x - min)
	// f(x) = --------------  + a
	// 		     max - min

    scaled := ((mend_max - mend_min) * float64(im_val) / float64(im_max)) + mend_min

    return scaled
}

// func scale_iteration_to_color(it_val int, it_max int, col_min int, col_max int) int{
// 	//        (b-a)(x - min)
// 	// f(x) = --------------  + a
// 	// 		     max - min

//     scaled := ((col_max - col_min) * (it_val - 1) / (it_max - 1)) + col_min

//     return scaled
// }

func mandelbrot(r_part float64, i_part float64, max_iteration int) (float64, int) {
	c := complex(r_part, i_part)
	var z complex128 = complex(0, 0)
	var mod float64
	n := 0

	for i := 0; i < max_iteration; i++ {
		n++
		z = cmplx.Pow(z,2) + c
		mod = math.Sqrt(math.Pow(real(z),2) + math.Pow(imag(z),2))
		if(mod > 2){
			break
		}
	}
	return mod, n
}

/*
Position = 0.0     Color = (  0,   7, 100)
Position = 0.16    Color = ( 32, 107, 203)
Position = 0.42    Color = (237, 255, 255)
Position = 0.6425  Color = (255, 170,   0)
Position = 0.8575  Color = (  0,   2,   0)
*/
func compute_colormap(n_color int) []color.NRGBA{
	
	xs := []float64{0.0, 0.16, 0.42, 0.6425, 0.8575, 1.0}
	yrs := []float64{0, 32, 237, 255, 0, 0}
	ygs := []float64{7, 107, 255, 170, 2, 7}
	ybs := []float64{100, 203, 255, 0, 0, 100}

	colormap := make([]color.NRGBA,0,n_color)

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
	}
	log.Print(colormap)
	return colormap
}

func get_color(it int, c complex128, colors []color.NRGBA) color.NRGBA {
	//found on https://stackoverflow.com/questions/16500656/which-color-gradient-is-used-to-color-mandelbrot-in-wikipedia
	smoothed := math.Log2(math.Log2(math.Pow(real(c),2) + math.Pow(imag(c),2)) / 2);
	colorI := (int)(math.Sqrt(float64(it) + 10.0 - smoothed) * 256) % len(colors);
	var color color.NRGBA
	if colorI == -8{
		//log.Print(c)
	} else {
		color = colors[colorI]
	}
	return color
}

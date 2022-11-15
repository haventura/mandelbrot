package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func main(){
	
	min_x := -2.0
	max_x := 1.0
	step_x := 0.001
	min_y := -1.0
	max_y := 1.0
	step_y := 0.001

	n := int(math.Ceil(((max_x - min_x) / step_x) * ((max_y - min_y) / step_y)))
	pts := make(plotter.XYs, 0, n)

	fmt.Printf("Computing...\n")
	start_comp := time.Now()
	for x := min_x; x < max_x; x += step_x  {
		for y := min_y; y < max_y; y += step_y {
			val, _ := mandelbrot(x,y,30)
			if(val < 2){
				pts = append(pts, plotter.XY{X: x, Y: y})
			}
			//fmt.Printf("(%v, %v), val: %v, it: %v\n", x, y, val, it)
		}
	}
	elapsed_comp := time.Since(start_comp)
	fmt.Printf("computing %v points took %s\n", n, elapsed_comp)
	fmt.Printf("Rendering...\n")
	start_rend := time.Now()
	p := plot.New()
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	p.Add(plotter.NewGrid())

	s, _ := plotter.NewScatter(pts)
	p.Add(s)
	p.Save(8000, 6000, "scatter.png")
	elapsed_rend := time.Since(start_rend)
	fmt.Printf("Rendering took %s", elapsed_rend)
}

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
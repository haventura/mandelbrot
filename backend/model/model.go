package model

type Image_data struct {
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Min_r         float64 `json:"min_r"`
	Max_r         float64 `json:"max_r"`
	Min_i         float64 `json:"min_i"`
	Max_i         float64 `json:"max_i"`
	Max_iteration int     `json:"max_iteration"`
	Colormap_name string  `json:"colormap_name"`
}

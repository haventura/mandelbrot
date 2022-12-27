package model

type Image_data struct {
	Width         int     `json:"width"`
	Height        int     `json:"height"`
	Max_iteration int     `json:"max_iteration"`
	Min_r         float64 `json:"min_r"`
	Max_r         float64 `json:"max_r"`
	Min_i         float64 `json:"min_i"`
	Max_i         float64 `json:"max_i"`
	Colormap_name string  `json:"colormap_name"`
}

type Image_chunk_data struct {
	Total_width   int     `json:"total_width"`
	Total_height  int     `json:"total_height"`
	Chunck_width  int     `json:"chunk_width"`
	Chunck_height int     `json:"chunk_height"`
	Chunck_min_x  int     `json:"chunk_min_x"`
	Chunck_min_y  int     `json:"chunk_min_y"`
	Chunck_min_r  float64 `json:"chunk_min_r"`
	Chunck_max_r  float64 `json:"chunk_max_r"`
	Chunck_min_i  float64 `json:"chunk_min_i"`
	Chunck_max_i  float64 `json:"chunk_max_i"`
	Max_iteration int     `json:"max_iteration"`
	Colormap_name string  `json:"colormap_name"`
}

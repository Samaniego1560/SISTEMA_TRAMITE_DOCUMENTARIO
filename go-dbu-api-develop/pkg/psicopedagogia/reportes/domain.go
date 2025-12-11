package reportes

import "time"

type HeaderAttentionExcel struct {
	Frame string `json:"frame"`
	Title string `json:"title"`
	Area  string `json:"area"`
}

type MergeRange struct {
	Column1 string
	Column2 string
}

type DateRange struct {
	Start time.Time
	End   time.Time
}

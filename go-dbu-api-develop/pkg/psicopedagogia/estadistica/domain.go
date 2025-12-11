package estadistica

type DailyStat struct {
	Date     string `json:"date" db:"date"`
	CountSrq int    `json:"countSrq" db:"countSrq"`
	Count    int    `json:"count" db:"count"`
}

type DataPie struct {
	Estados_evaluacion string `db:"estados_evaluacion" json:"estados_evaluacion"`
	Total              int    `db:"total" json:"total"`
}

type DataBarras struct {
	Escuela string `db:"escuela" json:"escuela"`
	Total   int    `db:"total" json:"total"`
}

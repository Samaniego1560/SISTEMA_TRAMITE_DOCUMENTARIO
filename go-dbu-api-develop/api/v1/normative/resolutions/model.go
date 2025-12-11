package resolutions

type ResponseResolution struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  any    `json:"data"`
	Code  int    `json:"code"`
	Type  string `json:"type"`
}

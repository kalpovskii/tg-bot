package req

type resModel struct {
	Result Result `json:"result"`
}

type Result struct {
	Price float64 `json:"price"`
	Error string  `json:"error"`
}

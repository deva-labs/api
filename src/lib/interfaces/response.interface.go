package interfaces

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Data   interface{} `json:"data"`
	Status Status      `json:"status"`
	Error  *string     `json:"error"`
}

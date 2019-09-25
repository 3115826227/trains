package response

type Response struct {
	Ret     int         `json:"ret"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

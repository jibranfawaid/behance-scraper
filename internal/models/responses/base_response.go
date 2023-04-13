package responses

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

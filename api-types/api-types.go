package apiTypes

type DefaultSuccessResponse struct {
	Message string `json:"message"`
}

type DefaultErrorResponse struct {
	Error string `json:"message"`
}

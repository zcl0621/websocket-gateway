package response

type ErrorResponse struct {
	Err string `json:"err"`
}
type ErrorResponseData struct {
	Err interface{} `json:"error_code"`
}
package dto

type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseMetadata struct {
	Status  int
	Code    string
	Message string
}

type BaseResponse struct {
	Data  interface{} `json:"data"`
	Error *ErrorData  `json:"error"`
}

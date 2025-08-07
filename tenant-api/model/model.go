package model

import "encoding/json"

const (
	// success already handle with response `success` true
	ErrorGeneral                  = "1001"
	ErrorInvalidRequest           = "1002" // wrong format field, ex: date etc.
	ErrorBadRequest               = "1003" // wrong field type, ex: json, text, etc.
	ErrorUnauthorized             = "1004"
	ErrorUnauthorizedTokenExpired = "1007"
	ErrorDataFound                = "1008"
	ErrorDuplicateData            = "1009"
)

type JsonResponse struct {
	Data       any     `json:"data,omitempty"`
	TotalData  *int    `json:"total_data,omitempty"`
	NextCursor *string `json:"next_cursor,omitempty"`
	Message    string  `json:"message,omitempty"`
	ErrorCode  string  `json:"error_code,omitempty"`
	Success    bool    `json:"success"`
}

func NewJsonResponse(success bool) *JsonResponse {
	return &JsonResponse{Success: success}
}

func NewError(code, message string) *JsonResponse {
	return &JsonResponse{Success: false, ErrorCode: code, Message: message}
}

func (r *JsonResponse) SetList(data any, total int) *JsonResponse {
	r.Data = data
	r.TotalData = &total
	return r
}

func (r *JsonResponse) SetData(data any) *JsonResponse {
	r.Data = data
	return r
}

func (r *JsonResponse) SetMessage(message string) *JsonResponse {
	r.Message = message
	return r
}

func (r *JsonResponse) SetListWithCursor(data any, cursor string) *JsonResponse {
	r.Data = data
	r.NextCursor = &cursor
	return r
}

func (r *JsonResponse) Error() string {
	errBytes, _ := json.Marshal(r)
	return string(errBytes)
}

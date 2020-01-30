package jumper

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	w http.ResponseWriter
	Status int `json:"status"`
	StatusNumber string `json:"status_number"`
	StatusCode string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Data interface{} `json:"data"`
}

func PlugResponse(w http.ResponseWriter) *Response {
	res := &Response{
		Status:        0,
		StatusNumber:  "",
		StatusCode:    "",
		StatusMessage: "",
		Data:          nil,
	}
	res.Assign(w)
	return res
}

func (r *Response) Assign(w http.ResponseWriter) {
	r.w = w
}

func (r *Response) SetHttpCode(code int) *Response {
	r.w.WriteHeader(code)
	return r
}

func (r *Response) Reply(status int, number string, code string, message string, data interface{}){
	r.w.Header().Set("Content-Type", "application/json")

	r.Status = status
	r.StatusNumber = number
	r.StatusCode = code
	r.StatusMessage = message
	r.Data = data

	err := json.NewEncoder(r.w).Encode(r)
	if err!=nil {
		log.Panic(err)
	}
}

func (r *Response) ReplyFailed(number string, code string, message string, data interface{}) {
	r.Reply(0, number, code, message, data)
}

func (r *Response) ReplySuccess(number string, code string, message string, data interface{}) {
	r.Reply(1, number, code, message, data)
}
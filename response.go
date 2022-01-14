package jumper

import (
	"encoding/json"
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
	res.w = w
	return res
}

func (r *Response) SetHttpCode(code int) *Response {
	r.w.WriteHeader(code)
	return r
}

// Reply 'data' arguments only used on index 0 */
func (r *Response) Reply(status int, number string, code string, message string, data... interface{}) error {
	r.w.Header().Set("Content-Type", "application/json")

	r.Status = status
	r.StatusNumber = number
	r.StatusCode = code
	r.StatusMessage = message
	if len(data) > 0 {
		r.Data = data[0]
	}

	return json.NewEncoder(r.w).Encode(r)
}

// ReplyFailed 'data' arguments only used on index 0 */
func (r *Response) ReplyFailed(number string, code string, message string, data... interface{}) error {
	return r.Reply(0, number, code, message, data...)
}

// ReplySuccess 'data' arguments only used on index 0 */
func (r *Response) ReplySuccess(number string, code string, message string, data... interface{}) error {
	return r.Reply(1, number, code, message, data...)
}
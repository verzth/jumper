package jumper

import (
	"encoding/json"
	"net/http"
)

type Response interface {
	SetHttpCode(code int) Response
	ReplyAs(res Response) error
	Reply(status int, number string, code string, message string, data ...interface{}) error
	ReplyFailed(number string, code string, message string, data ...interface{}) error
	ReplySuccess(number string, code string, message string, data ...interface{}) error
	HttpStatusCode() int
	SetHttpStatusCode(httpStatusCode int) Response
	GetStatus() int
	GetStatusNumber() string
	GetStatusCode() string
	GetStatusMessage() string
	GetData() interface{}
}

type ResponseX struct {
	w              http.ResponseWriter
	httpStatusCode int
	Status         int         `json:"status"`
	StatusNumber   string      `json:"status_number"`
	StatusCode     string      `json:"status_code"`
	StatusMessage  string      `json:"status_message"`
	Data           interface{} `json:"data"`
}

func NewResponse(httpStatusCode int, Status int, StatusNumber string, StatusCode string, StatusMessage string, Data interface{}) Response {
	rx := ResponseX{
		Status:        Status,
		StatusNumber:  StatusNumber,
		StatusCode:    StatusCode,
		StatusMessage: StatusMessage,
		Data:          Data,
	}
	rx.SetHttpCode(httpStatusCode)
	return &rx
}

func (r *ResponseX) HttpStatusCode() int {
	return r.httpStatusCode
}

func (r *ResponseX) SetHttpStatusCode(httpStatusCode int) Response {
	r.httpStatusCode = httpStatusCode
	return r
}

func (r *ResponseX) GetStatus() int {
	return r.Status
}

func (r *ResponseX) GetStatusNumber() string {
	return r.StatusNumber
}

func (r *ResponseX) GetStatusCode() string {
	return r.StatusCode
}

func (r *ResponseX) GetStatusMessage() string {
	return r.StatusMessage
}

func (r *ResponseX) GetData() interface{} {
	return r.Data
}

func PlugResponse(w http.ResponseWriter) Response {
	res := &ResponseX{
		Status:        0,
		StatusNumber:  "",
		StatusCode:    "",
		StatusMessage: "",
		Data:          nil,
	}
	res.w = w
	return res
}

func (r *ResponseX) SetHttpCode(code int) Response {
	r.w.WriteHeader(code)
	return r
}

func (r *ResponseX) ReplyAs(res Response) error {
	if res.HttpStatusCode() != 0 {
		r.w.WriteHeader(res.HttpStatusCode())
	}
	return r.Reply(res.GetStatus(), res.GetStatusNumber(), res.GetStatusCode(), res.GetStatusMessage(), res.GetData())
}

// Reply 'data' arguments only used on index 0 */
func (r *ResponseX) Reply(status int, number string, code string, message string, data ...interface{}) error {
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
func (r *ResponseX) ReplyFailed(number string, code string, message string, data ...interface{}) error {
	return r.Reply(0, number, code, message, data...)
}

// ReplySuccess 'data' arguments only used on index 0 */
func (r *ResponseX) ReplySuccess(number string, code string, message string, data ...interface{}) error {
	return r.Reply(1, number, code, message, data...)
}

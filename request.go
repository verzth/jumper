package jumper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Request struct {
	r http.Request
	segments map[string]string
	params Params
	files map[string]interface{}
	header http.Header
	Method string
	ClientIP string
	ClientPort string
}

func PlugRequest(r *http.Request, w http.ResponseWriter) *Request {
	req := &Request{
		r:          *r,
		segments: mux.Vars(r),
		params: Params{},
		files: map[string]interface{}{},
		header: r.Header,
		Method: r.Method,
		ClientIP: getHost(r),
		ClientPort: getPort(r),
	}

	// PARSE QUERY STRING PARAMETERS
	for k, v := range r.URL.Query() {
		req.params[k] = scan(v)
	}

	switch r.Method {
	case http.MethodPut,http.MethodPost,http.MethodDelete:{
		contentType := req.header.Get("Content-Type")
		if strings.Contains(contentType, "multipart/form-data") {
			contentType = "multipart/form-data"
		}
		switch contentType {
		case "multipart/form-data":{
			err := r.ParseMultipartForm(32 << 10)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return req
			}
			for k, v := range r.MultipartForm.Value {
				req.params[k] = scan(v)
			}
			for k, v := range r.MultipartForm.File {
				req.files[k] = scanFiles(v)
			}
			break
		}
		case "application/x-www-form-urlencoded":{
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return req
			}
			for k, v := range r.PostForm {
				req.params[k] = scan(v)
			}
			break
		}
		case "application/json":{
			dec := json.NewDecoder(r.Body)

			err := dec.Decode(&req.params)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return req
			}
			break
		}
		}
		break
	}
	}
	return req
}

func scan(values []string) interface{} {
	if len(values) == 1 {
		return values[0]
	}else if len(values) > 1 {
		return values
	} else {
		return nil
	}
}

func scanFiles(values []*multipart.FileHeader) interface{} {
	if len(values) == 1 {
		return values[0]
	}else if len(values) > 1 {
		return values
	} else {
		return nil
	}
}


func (r *Request) GetHost() string {
	return r.r.URL.Hostname()
}

func (r *Request) GetPort() string {
	return r.r.URL.Port()
}

func (r *Request) GetScheme() string {
	return r.r.URL.Scheme
}

func (r *Request) GetOpaque() string {
	return r.r.URL.Opaque
}

func (r *Request) GetPath() string {
	return r.r.URL.Path
}

func (r *Request) GetRawPath() string {
	return r.r.URL.RawPath
}

func (r *Request) GetRawQuery() string {
	return r.r.URL.RawQuery
}

func (r *Request) GetFragment() string {
	return r.r.URL.Fragment
}

func (r *Request) HasUser() bool {
	_,_,ok := r.r.BasicAuth()
	return ok
}

func (r *Request) GetUsername() string {
	user,_,ok := r.r.BasicAuth()
	if ok {
		return user
	}
	return ""
}

func (r *Request) GetPassword() string {
	_,pass,ok := r.r.BasicAuth()
	if ok {
		return pass
	}
	return ""
}

func (r *Request) GetUrl() string {
	return r.r.URL.Scheme+"://"+r.r.URL.Host+r.r.URL.EscapedPath()
}

func (r *Request) GetFullUrl() string {
	return r.r.URL.String()
}

func (r *Request) Header(key string) string {
	return r.header.Get(key)
}

func (r *Request) GetAll() map[string] interface{} {
	return r.params
}

func (r *Request) Get(key string) string {
	if r.params[key] != nil {
		return fmt.Sprintf("%v", r.params[key])
	}
	return ""
}

func (r *Request) Append(key string, val string) {
	r.params[key] = val
}

func (r *Request) GetSegment(key string) string {
	return r.segments[key]
}

func (r *Request) GetFile(key string) (*File, error) {
	if r.files[key] != nil {
		_, ok := r.files[key].(*multipart.FileHeader)
		if ok {
			f, err := r.files[key].(*multipart.FileHeader).Open()
			return &File{
				f:  f,
				fh: r.files[key].(*multipart.FileHeader),
			}, err
		} else {
			return nil, errors.New("invalid file, maybe files instead")
		}
	}
	return nil, errors.New("no such file")
}

func (r *Request) GetFiles(key string) ([]*File, error) {
	if r.files[key] != nil {
		var files []*File
		vs, ok := r.files[key].([]*multipart.FileHeader)
		if ok {
			for _, v := range vs{
				f, err := v.Open()
				if err != nil {
					return nil, errors.New("files error")
				}
				files = append(files, &File{
					f:  f,
					fh: v,
				})
			}
			return files, nil
		} else {
			return nil, errors.New("invalid files, maybe file instead")
		}
	}
	return nil, errors.New("no such file")
}

func (r *Request) GetUint64(key string) uint64 {
	if r.params[key] != nil {
		switch r.params[key].(type) {
		case float64: return uint64(r.params[key].(float64))
		case int: return uint64(r.params[key].(int))
		case string:
			i64, _ := strconv.ParseUint(r.params[key].(string), 10, 32)
			return i64
		}
	}
	return 0
}

func (r *Request) GetUint32(key string) uint32 {
	return uint32(r.GetUint64(key))
}

func (r *Request) GetUint(key string) uint {
	return uint(r.GetUint64(key))
}

func (r *Request) GetInt64(key string) int64 {
	if r.params[key] != nil {
		switch r.params[key].(type) {
		case float64: return int64(r.params[key].(float64))
		case int: return int64(r.params[key].(int))
		case string:
			i64, _ := strconv.ParseInt(r.params[key].(string), 10, 32)
			return i64
		}
	}
	return 0
}

func (r *Request) GetInt32(key string) int32 {
	return int32(r.GetInt64(key))
}

func (r *Request) GetInt(key string) int {
	return int(r.GetInt64(key))
}

func (r *Request) GetFloat64(key string) float64 {
	if r.params[key] != nil {
		switch r.params[key].(type) {
		case float64: return r.params[key].(float64)
		case int: return float64(r.params[key].(int))
		case string:
			i64, _ := strconv.ParseFloat(r.params[key].(string), 10)
			return i64
		}
	}
	return 0
}

func (r *Request) GetFloat(key string) float32 {
	return float32(r.GetFloat64(key))
}

func (r *Request) GetTime(key string) (*time.Time,error) {
	if r.params[key] != nil {
		t, err := time.Parse(time.RFC3339,r.params[key].(string))
		if err != nil {
			return nil, errors.New("use RFC3339 format string for datetime")
		}
		return &t, nil
	} else {
		return nil, errors.New("no time specified")
	}
}

func (r *Request) GetTimeNE(key string) *time.Time {
	t, _ := r.GetTime(key)
	return t
}

func (r *Request) GetArray(key string) []interface{} {
	if r.params[key] != nil {
		if v, ok := r.params[key].([]interface{}); ok {
			return v
		}
	}
	return nil
}

func (r *Request) GetMap(key string) map[string]interface{} {
	if r.params[key] != nil {
		if v, ok := r.params[key].(map[string]interface{}); ok {
			return v
		}
	}
	return nil
}

func (r *Request) GetJSON(key string) JSON {
	jsonObj, err := json.Marshal(r.params[key])
	if err != nil {
		return nil
	}else{
		return jsonObj
	}
}

func (r *Request) GetStruct(obj interface{}) error {
	decoder := json.NewDecoder(r.r.Body)
	return decoder.Decode(&obj)
}

func (r *Request) has(key string) bool {
	if _, found := r.params[key]; !found {
		return false
	}
	return true
}

func (r *Request) Has(keys... string) bool {
	found := true
	for _, key := range keys {
		found = found && r.has(key)
	}
	return found
}

func (r *Request) Filled(keys... string) bool {
	found := true
	for _, key := range keys {
		found = found && r.Has(key)
		if _, ok := r.params[key].(string); ok {
			found = found && strings.TrimSpace(r.params[key].(string)) != ""
		}
	}
	return found
}

func (r *Request) HasFile(keys... string) bool {
	found := true
	for _, key := range keys {
		found = found && r.files[key] != nil
	}
	return found
}
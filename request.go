package jumper

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"git.verzth.work/go/utils"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type RequestJumper struct {
	r          http.Request
	segments   map[string]string
	params     Params
	files      map[string]interface{}
	header     http.Header
	Method     string
	ClientIP   string
	ClientPort string
}

func PlugRequest(r *http.Request, w http.ResponseWriter) *RequestJumper {
	req := &RequestJumper{
		r:          *r,
		segments:   mux.Vars(r),
		params:     Params{},
		files:      map[string]interface{}{},
		header:     r.Header,
		Method:     r.Method,
		ClientIP:   getHost(r),
		ClientPort: getPort(r),
	}

	// PARSE QUERY STRING PARAMETERS
	for k, v := range r.URL.Query() {
		req.params[k] = scan(v)
	}

	switch r.Method {
	case http.MethodGet, "FETCH", http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch:
		{
			contentType := req.header.Get("Content-Type")
			if strings.Contains(contentType, "multipart/form-data") {
				if r.Method == http.MethodGet {
					return req
				}
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
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				if r.Method == http.MethodGet {
					return req
				}
				err := r.ParseForm()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return req
				}
				for k, v := range r.PostForm {
					req.params[k] = scan(v)
				}
			} else if strings.Contains(contentType, "application/json") {
				if r.ContentLength > 0 {
					dec := json.NewDecoder(r.Body)

					err := dec.Decode(&req.params)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return req
					}
				}
			}
			break
		}
	}
	return req
}

// TouchRequest touch request with rewrite to reader, so handler can reuse the reader.
func TouchRequest(r *http.Request, w http.ResponseWriter) *RequestJumper {
	req := &RequestJumper{
		r:          *r,
		segments:   mux.Vars(r),
		params:     Params{},
		files:      map[string]interface{}{},
		header:     r.Header,
		Method:     r.Method,
		ClientIP:   getHost(r),
		ClientPort: getPort(r),
	}

	// PARSE QUERY STRING PARAMETERS
	for k, v := range r.URL.Query() {
		req.params[k] = scan(v)
	}

	switch r.Method {
	case http.MethodGet, "FETCH", http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch:
		{
			contentType := req.header.Get("Content-Type")
			if strings.Contains(contentType, "multipart/form-data") {
				if r.Method == http.MethodGet {
					return req
				}
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
			} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
				if r.Method == http.MethodGet {
					return req
				}
				err := r.ParseForm()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return req
				}
				for k, v := range r.PostForm {
					req.params[k] = scan(v)
				}
			} else if strings.Contains(contentType, "application/json") {
				if r.ContentLength > 0 {
					b := bytes.NewBuffer(make([]byte, 0))
					reader := io.TeeReader(r.Body, b)

					dec := json.NewDecoder(reader)

					err := dec.Decode(&req.params)

					r.Body = ioutil.NopCloser(b)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return req
					}
				}
			}
			break
		}
	}
	return req
}

func scan(values []string) interface{} {
	if len(values) == 1 {
		return identify(values[0])
	} else if len(values) > 1 {
		list := []interface{}{}
		for k, vs := range values {
			list[k] = identify(vs)
		}
		return list
	} else {
		return nil
	}
}

func identify(value string) interface{} {
	var arr []interface{}
	var mp map[string]interface{}
	errArr := json.Unmarshal([]byte(value), &arr)
	errMp := json.Unmarshal([]byte(value), &mp)
	if errArr == nil {
		return arr
	} else if errMp == nil {
		return mp
	} else {
		return value
	}
}

func scanFiles(values []*multipart.FileHeader) interface{} {
	if len(values) == 1 {
		return values[0]
	} else if len(values) > 1 {
		return values
	} else {
		return nil
	}
}

func (r *RequestJumper) GetHost() string {
	return r.r.URL.Hostname()
}

func (r *RequestJumper) GetPort() string {
	return r.r.URL.Port()
}

func (r *RequestJumper) GetScheme() string {
	return r.r.URL.Scheme
}

func (r *RequestJumper) GetOpaque() string {
	return r.r.URL.Opaque
}

func (r *RequestJumper) GetPath() string {
	return r.r.URL.Path
}

func (r *RequestJumper) GetRawPath() string {
	return r.r.URL.RawPath
}

func (r *RequestJumper) GetRawQuery() string {
	return r.r.URL.RawQuery
}

func (r *RequestJumper) GetFragment() string {
	return r.r.URL.Fragment
}

func (r *RequestJumper) HasUser() bool {
	_, _, ok := r.r.BasicAuth()
	return ok
}

func (r *RequestJumper) GetUsername() string {
	user, _, ok := r.r.BasicAuth()
	if ok {
		return user
	}
	return ""
}

func (r *RequestJumper) GetPassword() string {
	_, pass, ok := r.r.BasicAuth()
	if ok {
		return pass
	}
	return ""
}

func (r *RequestJumper) GetUrl() string {
	return r.r.URL.Scheme + "://" + r.r.URL.Host + r.r.URL.EscapedPath()
}

func (r *RequestJumper) GetFullUrl() string {
	return r.r.URL.String()
}

func (r *RequestJumper) Header(key string) string {
	return r.header.Get(key)
}

func (r *RequestJumper) Append(key string, val string) {
	r.params[key] = val
}

func (r *RequestJumper) GetSegment(key string) string {
	return r.segments[key]
}

func (r *RequestJumper) GetSegmentUint64(key string) uint64 {
	if r.segments[key] != "" {
		i64, _ := strconv.ParseUint(r.segments[key], 10, 32)
		return i64
	}
	return 0
}

func (r *RequestJumper) GetSegmentUint32(key string) uint32 {
	return uint32(r.GetSegmentUint64(key))
}

func (r *RequestJumper) GetSegmentUint(key string) uint {
	return uint(r.GetSegmentUint64(key))
}

func (r *RequestJumper) GetSegmentInt64(key string) int64 {
	if r.segments[key] != "" {
		i64, _ := strconv.ParseInt(r.segments[key], 10, 32)
		return i64
	}
	return 0
}

func (r *RequestJumper) GetSegmentInt32(key string) int32 {
	return int32(r.GetSegmentInt64(key))
}

func (r *RequestJumper) GetSegmentInt(key string) int {
	return int(r.GetSegmentInt64(key))
}

func (r *RequestJumper) GetFile(key string) (*File, error) {
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

func (r *RequestJumper) GetFiles(key string) ([]*File, error) {
	if r.files[key] != nil {
		var files []*File
		vs, ok := r.files[key].([]*multipart.FileHeader)
		if ok {
			for _, v := range vs {
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

func (r *RequestJumper) GetAll() map[string]interface{} {
	return r.params
}

func ParseTo[T any](r *RequestJumper) (T, error) {
	jsonString, _ := json.Marshal(r.params)
	var en T
	err := json.Unmarshal(jsonString, &en)
	if err != nil {
		return nil, err
	} else {
		return en, nil
	}
}

func ParseOf[T any](r *RequestJumper, en *T) error {
	jsonString, _ := json.Marshal(r.params)
	err := json.Unmarshal(jsonString, en)
	return err
}

func (r *RequestJumper) GetPtr(key string) *interface{} {
	val := reflect.ValueOf(r.params[key])
	if r.params[key] != nil || (val.IsValid() && val.Kind() == reflect.Interface) {
		v := r.params[key]
		return &v
	}
	return nil
}

func (r *RequestJumper) Get(key string) interface{} {
	v := r.GetPtr(key)
	if v != nil {
		return v
	} else {
		return nil
	}
}

func (r *RequestJumper) GetStringPtr(key string) *string {
	val := reflect.ValueOf(r.params[key])
	if r.params[key] != nil || (val.IsValid() && val.Kind() == reflect.Slice && val.Len() > 0) {
		v := fmt.Sprintf("%v", r.params[key])
		return &v
	}
	return nil
}

func (r *RequestJumper) GetString(key string) string {
	v := r.GetStringPtr(key)
	if v != nil {
		return *v
	} else {
		return ""
	}
}

func (r *RequestJumper) GetUint64Ptr(key string) *uint64 {
	if r.params[key] != nil {
		var v uint64
		switch r.params[key].(type) {
		case float64:
			v = uint64(r.params[key].(float64))
		case int:
			v = uint64(r.params[key].(int))
		case string:
			v, _ = strconv.ParseUint(r.params[key].(string), 10, 32)
		case bool:
			{
				if r.params[key].(bool) {
					v = 1
				} else {
					v = 0
				}
			}
		}
		return &v
	}
	return nil
}

func (r *RequestJumper) GetUint64(key string) uint64 {
	v := r.GetUint64Ptr(key)
	if v != nil {
		return *v
	} else {
		return 0
	}
}

func (r *RequestJumper) GetUint32Ptr(key string) *uint32 {
	v := r.GetUint64Ptr(key)
	if v != nil {
		val := uint32(*v)
		return &val
	} else {
		return nil
	}
}

func (r *RequestJumper) GetUint32(key string) uint32 {
	return uint32(r.GetUint64(key))
}

func (r *RequestJumper) GetUintPtr(key string) *uint {
	v := r.GetUint64Ptr(key)
	if v != nil {
		val := uint(*v)
		return &val
	} else {
		return nil
	}
}

func (r *RequestJumper) GetUint(key string) uint {
	return uint(r.GetUint64(key))
}

func (r *RequestJumper) GetInt64Ptr(key string) *int64 {
	if r.params[key] != nil {
		var v int64
		switch r.params[key].(type) {
		case float64:
			v = int64(r.params[key].(float64))
		case int:
			v = int64(r.params[key].(int))
		case string:
			v, _ = strconv.ParseInt(r.params[key].(string), 10, 32)
		case bool:
			{
				if r.params[key].(bool) {
					v = 1
				} else {
					v = 0
				}
			}
		}
		return &v
	}
	return nil
}

func (r *RequestJumper) GetInt64(key string) int64 {
	v := r.GetInt64Ptr(key)
	if v != nil {
		return *v
	} else {
		return 0
	}
}

func (r *RequestJumper) GetInt32Ptr(key string) *int32 {
	v := r.GetUint64Ptr(key)
	if v != nil {
		val := int32(*v)
		return &val
	} else {
		return nil
	}
}

func (r *RequestJumper) GetInt32(key string) int32 {
	return int32(r.GetInt64(key))
}

func (r *RequestJumper) GetIntPtr(key string) *int {
	v := r.GetUint64Ptr(key)
	if v != nil {
		val := int(*v)
		return &val
	} else {
		return nil
	}
}

func (r *RequestJumper) GetInt(key string) int {
	return int(r.GetInt64(key))
}

func (r *RequestJumper) GetFloat64Ptr(key string) *float64 {
	if r.params[key] != nil {
		var v float64
		switch r.params[key].(type) {
		case float64:
			v = r.params[key].(float64)
		case int:
			v = float64(r.params[key].(int))
		case string:
			v, _ = strconv.ParseFloat(r.params[key].(string), 10)
		case bool:
			{
				if r.params[key].(bool) {
					v = 1
				} else {
					v = 0
				}
			}
		}
		return &v
	}
	return nil
}

func (r *RequestJumper) GetFloat64(key string) float64 {
	v := r.GetFloat64Ptr(key)
	if v != nil {
		return *v
	} else {
		return 0
	}
}

func (r *RequestJumper) GetFloat32Ptr(key string) *float32 {
	v := r.GetFloat64Ptr(key)
	if v != nil {
		val := float32(*v)
		return &val
	} else {
		return nil
	}
}

func (r *RequestJumper) GetFloat(key string) float32 {
	return float32(r.GetFloat64(key))
}

func (r *RequestJumper) GetBoolPtr(key string) *bool {
	if r.params[key] != nil {
		var v bool
		switch r.params[key].(type) {
		case float64:
			v = r.params[key].(float64) > 0
		case int:
			v = float64(r.params[key].(int)) > 0
		case string:
			i64, _ := strconv.ParseFloat(r.params[key].(string), 10)
			v = i64 > 0
		case bool:
			v = r.params[key].(bool)
		}
		return &v
	}
	return nil
}

func (r *RequestJumper) GetBool(key string) bool {
	v := r.GetBoolPtr(key)
	if v != nil {
		return *v
	} else {
		return false
	}
}

func (r *RequestJumper) GetTime(key string) (*time.Time, error) {
	if r.params[key] != nil {
		t, err := time.Parse(time.RFC3339, r.params[key].(string))
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05.000Z07:00", r.params[key].(string)) // RFC3339Mili
			if err != nil {
				t, err = time.Parse(time.RFC3339Nano, r.params[key].(string))
				if err != nil {
					return nil, errors.New("use RFC3339 format string for datetime")
				}
			}
		}
		return &t, nil
	} else {
		return nil, errors.New("no time specified")
	}
}

func (r *RequestJumper) GetTimeNE(key string) *time.Time {
	t, _ := r.GetTime(key)
	return t
}

func (r *RequestJumper) GetArray(key string) []interface{} {
	if r.params[key] != nil {
		if v, ok := r.params[key].([]interface{}); ok {
			return v
		}
	}
	return nil
}

func (r *RequestJumper) GetArrayUniquify(key string) []interface{} {
	if r.params[key] != nil {
		if v, ok := r.params[key].([]interface{}); ok {
			utils.Slice.Uniquify(&v)
			return v
		}
	}
	return nil
}

func (r *RequestJumper) GetMap(key string) map[string]interface{} {
	if r.params[key] != nil {
		if v, ok := r.params[key].(map[string]interface{}); ok {
			return v
		}
	}
	return nil
}

func (r *RequestJumper) GetJSON(key string) JSON {
	jsonObj, err := json.Marshal(r.params[key])
	if err != nil {
		return nil
	} else {
		return jsonObj
	}
}

func (r *RequestJumper) GetStruct(obj interface{}) error {
	decoder := json.NewDecoder(r.r.Body)
	return decoder.Decode(&obj)
}

func (r *RequestJumper) has(key string) bool {
	if _, found := r.params[key]; !found {
		return false
	}
	return true
}

func (r *RequestJumper) Has(keys ...string) (found bool) {
	found = true
	for _, key := range keys {
		found = found && r.has(key)
	}
	return
}

func (r *RequestJumper) Filled(keys ...string) (found bool) {
	found = true
	for _, key := range keys {
		found = found && r.has(key)
		val := reflect.ValueOf(r.params[key])
		if val.IsValid() {
			switch val.Kind() {
			case reflect.String:
				found = found && strings.TrimSpace(r.GetString(key)) != ""
			case reflect.Slice:
				found = found && val.Len() > 0
			case reflect.Array:
				found = found && val.Len() > 0
			}
		} else {
			found = false
		}
	}
	return
}

func (r *RequestJumper) hasHeader(key string) bool {
	if _, found := r.header[textproto.CanonicalMIMEHeaderKey(key)]; !found {
		return false
	}
	return true
}

func (r *RequestJumper) HasHeader(keys ...string) (found bool) {
	found = true
	for _, key := range keys {
		found = found && r.hasHeader(key)
	}
	return
}

func (r *RequestJumper) HeaderFilled(keys ...string) (found bool) {
	found = true
	for _, key := range keys {
		found = found && r.hasHeader(key) && r.Header(key) != ""
	}
	return
}

func (r *RequestJumper) HasFile(keys ...string) (found bool) {
	found = true
	for _, key := range keys {
		found = found && r.files[key] != nil
	}
	return
}

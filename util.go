package jumper

import (
	"net/http"
	"strings"
)

type Params map[string]interface{}

func getHost(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	lIndex := strings.LastIndex(IPAddress,":")
	if lIndex != -1 {
		return string([]rune(IPAddress)[0:lIndex])
	}else {
		return string([]rune(IPAddress)[0:])
	}
}

func getPort(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	lIndex := strings.LastIndex(IPAddress,":")
	if lIndex != -1 {
		return string([]rune(IPAddress)[lIndex+1:])
	}else {
		return ""
	}
}
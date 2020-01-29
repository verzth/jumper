package jumper

import (
	"net/http"
	"strings"
)

func getHost(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	cs := strings.Split(IPAddress, ":")
	return cs[0]
}

func getPort(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	cs := strings.Split(r.RemoteAddr, ":")

	if len(cs) == 2 {
		return cs[1]
	}else{
		return ""
	}
}
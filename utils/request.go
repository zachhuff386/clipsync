package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func StripPort(hostport string) string {
	colon := strings.IndexByte(hostport, ':')
	if colon == -1 {
		return hostport
	}
	if i := strings.IndexByte(hostport, ']'); i != -1 {
		return strings.TrimPrefix(hostport[:i], "[")
	}
	return hostport[:colon]
}

func FormatHostPort(hostname string, port int) string {
	if strings.Contains(hostname, ":") {
		hostname = "[" + hostname + "]"
	}
	return fmt.Sprintf("%s:%d", hostname, port)
}

func GetStatusMessage(code int) string {
	return fmt.Sprintf("%d %s", code, http.StatusText(code))
}

func WriteStatus(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, GetStatusMessage(code))
}

func WriteText(w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, text)
}

func WriteUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(401)
	fmt.Fprintln(w, "401 "+msg)
}

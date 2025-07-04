package http

import (
	"encoding/json"
	"net/http"

	"github.com/Nemutagk/goerrors"
)

var validStatusCodes = map[int]struct{}{
	100: {}, 101: {}, 102: {}, 103: {},
	200: {}, 201: {}, 202: {}, 203: {}, 204: {}, 205: {}, 206: {}, 207: {}, 208: {}, 226: {},
	300: {}, 301: {}, 302: {}, 303: {}, 304: {}, 305: {}, 307: {}, 308: {},
	400: {}, 401: {}, 402: {}, 403: {}, 404: {}, 405: {}, 406: {}, 407: {}, 408: {}, 409: {},
	410: {}, 411: {}, 412: {}, 413: {}, 414: {}, 415: {}, 416: {}, 417: {}, 418: {}, 421: {},
	422: {}, 423: {}, 424: {}, 425: {}, 426: {}, 428: {}, 429: {}, 431: {}, 451: {},
	500: {}, 501: {}, 502: {}, 503: {}, 504: {}, 505: {}, 506: {}, 507: {}, 508: {}, 510: {}, 511: {},
}

func Success(w http.ResponseWriter, data any, code int, contentType string, headers ...map[string]string) {
	w.Header().Set("Content-Type", contentType)

	for _, header := range headers {
		for key, value := range header {
			if key == "Content-Type" && value != "" {
				w.Header().Set(key, value)
			}
		}
	}

	if _, valid := validStatusCodes[code]; !valid {
		code = http.StatusInternalServerError
	}

	w.WriteHeader(code)

	// Detección automática del tipo de dato a enviar
	if contains(contentType, []string{"json", "text", "xml"}) {
		// Si es JSON, intentamos codificar como JSON
		if contains(contentType, []string{"json"}) {
			if err := json.NewEncoder(w).Encode(data); err != nil {
				http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			}
			return
		}
		// Si es texto o xml, lo convertimos a string
		str, ok := data.(string)
		if !ok {
			http.Error(w, "Data is not a string", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(str))
		return
	}

	if contains(contentType, []string{"image", "audio", "video", "application"}) {
		bytes, ok := data.([]byte)
		if !ok {
			http.Error(w, "Data is not a byte slice", http.StatusInternalServerError)
			return
		}
		w.Write(bytes)
		return
	}

	http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
}

func Error(w http.ResponseWriter, err goerrors.GError) {
	Success(w, err, err.Status, "application/json")
}

// contains verifica si el contentType contiene alguno de los substrings dados
func contains(contentType string, substrings []string) bool {
	for _, sub := range substrings {
		if len(contentType) >= len(sub) && (contentType == sub || (len(contentType) > len(sub) && containsStr(contentType, sub))) {
			return true
		}
	}
	return false
}

// containsStr es un helper para saber si un substring está en un string (case-insensitive)
func containsStr(s, substr string) bool {
	return (len(s) >= len(substr)) && (stringIndexInsensitive(s, substr) >= 0)
}

// stringIndexInsensitive busca un substring sin importar mayúsculas/minúsculas
func stringIndexInsensitive(s, substr string) int {
	return indexOfInsensitive(s, substr)
}

// indexOfInsensitive busca la posición de un substring sin importar mayúsculas/minúsculas
func indexOfInsensitive(s, substr string) int {
	return indexOf([]rune(lower(s)), []rune(lower(substr)))
}

func lower(s string) string {
	b := []rune(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func indexOf(s, substr []rune) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return i
		}
	}
	return -1
}

package goerrors

import "encoding/json"

type GError struct {
	Success  bool        `json:"success"`
	Message  string      `json:"message"`
	Errors   *[]SubError `json:"errors,omitempty"`
	Previous *GError     `json:"previous,omitempty"` // Optional field for chaining errors
	Status   int         `json:"status"`
}

type SubError struct {
	Type    string  `json:"type"`
	Message string  `json:"message"`
	Field   *string `json:"field"`
	Code    *string `json:"code"`
}

func NewGError(success bool, message string, status int, errors *[]SubError, previous *GError) *GError {
	return &GError{
		Success:  success,
		Message:  message,
		Status:   status,
		Errors:   errors,
		Previous: previous,
	}
}

func (e *GError) Error() string {
	return e.Message
}

func (e *GError) ToJson() string {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return `{"success": false, "message": "Error converting to JSON", "status": 500}`
	}
	return string(jsonData)
}

func (e *GError) ToMap() map[string]interface{} {
	var previous interface{}

	if e.Previous != nil {
		previous = e.Previous.ToMap()
	} else {
		previous = nil
	}

	return map[string]interface{}{
		"success":  e.Success,
		"message":  e.Message,
		"errors":   e.Errors,
		"status":   e.Status,
		"previous": previous,
	}
}

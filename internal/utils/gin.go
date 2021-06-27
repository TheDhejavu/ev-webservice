package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(err error) gin.H {
	return gin.H{
		"message": err.Error(),
	}
}

func FormatError(err map[string]string) gin.H {
	var _error = make(map[string]string)
	for k, v := range err {
		newKey := strings.Split(k, ".")
		key := newKey[1]
		key = strings.ToLower(key)
		_error[key] = v
	}
	return gin.H{
		"data":    _error,
		"message": "An error occurred",
	}
}

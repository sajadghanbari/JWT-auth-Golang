package helper

import (
	"JWT-Authentication-go/pkg/service_errors"
	"net/http"
)

var StatusCodeMapping = map[string]int{
	service_errors.OptExists:   409,
	service_errors.OtpUsed:     409,
	service_errors.OtpInvalid: 400,
}


func TranslateErrorToStatusCode(err error) int {
	value , ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
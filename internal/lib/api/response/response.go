package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strconv"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Alias  string `json:"alias,omitempty"`
}

const (
	StatusOK    = 200
	StatusError = "500"
)

func OK() Response {
	return Response{
		Status: strconv.Itoa(StatusOK),
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a required failed"), err.Field())
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is a required failed"), err.Field())
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("failed %s is not valid"), err.Field())
		}
	}
	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, " "),
	}
}

package helper

import (
	"github.com/gin-gonic/gin"
)

func SuccessResponse(status string, extendedData interface{}) gin.H {

	data := gin.H{
		"MESSAGE": "SUCCESS",
		"STATUS":  status,
	}

	if extendedData != nil {
		data["DATA"] = extendedData
	}

	return data
}

func ErrorResponse(status string, err string) gin.H {

	data := gin.H{
		"STATUS": status,
	}

	if err != "" {
		data["ERROR"] = err
	}

	return data
}

func BadRequestResponse(err string) gin.H {

	data := gin.H{
		"STATUS": "BAD_REQUEST",
	}

	if err != "" {
		data["ERROR"] = err
	}

	return data
}

func UnauthorizedResponse() gin.H {

	data := gin.H{
		"STATUS": "UNAUTHORIZED",
	}

	return data
}

func InternalErrorResponse(err string) gin.H {

	data := gin.H{
		"STATUS": "INTERNAL_SERVER_ERROR",
	}

	if err != "" {
		data["ERROR"] = err
	}

	return data
}

func NotFoundResponse(message string) gin.H {

	data := gin.H{
		"STATUS": "NOT_FOUND",
	}

	if message != "" {
		data["MESSAGE"] = message
	}

	return data
}

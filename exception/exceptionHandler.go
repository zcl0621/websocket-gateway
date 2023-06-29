package exception

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ws-gateway/logger"
)

func ExceptionHandler(context *gin.Context) {
	defer func() {
		var err = recover()
		if err != nil {
			processError(err, context)
			context.Abort()
		}
	}()
	context.Next()
}

func processError(exception interface{}, context *gin.Context) {
	switch t := exception.(type) {
	case ErrorContextProcess:
		t.printError()
		t.buildContext(context)
	default:
		logger.Logger("errorLogPrint", "error", exception.(error), nil)
		context.Writer.Header().Set("Error-Code", "1")
		context.JSON(http.StatusBadRequest, "操作失败")
	}
}

type ErrorContextProcess interface {
	printError()
	buildContext(context *gin.Context)
}

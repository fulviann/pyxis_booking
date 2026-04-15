package respond

import (
	"fmt"

	apierror "github.com/devanadindra/signlink-mobile/back-end/utils/api-error"
	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, err error) {

	apiErrors := apierror.GetApiErrors(err)
	ctx.JSON(apiErrors.Code, ApiModel[*string]{
		Data:   nil,
		Errors: apiErrors.Messages,
	})

	ctx.Set("error", err)
	ctx.Abort()
}

func Success(ctx *gin.Context, code int, data any) {
	ctx.Set("error", nil)
	if data == nil {
		ctx.Status(code)
		return
	}
	ctx.JSON(code, ApiModel[any]{
		Data:   data,
		Errors: nil,
	})
	ctx.Abort()
}

func Data(ctx *gin.Context, param DataParam) {
	ctx.Set("error", nil)
	if param.Data == nil {
		ctx.Status(param.Code)
		return
	}

	// Set the response headers
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", param.Filename))
	ctx.Header("Content-Type", param.MimeType)

	// Write the JSON data to the response
	ctx.Data(param.Code, param.MimeType, param.Data)

	ctx.Abort()
}

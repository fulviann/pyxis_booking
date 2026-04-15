package context

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apierror "github.com/devanadindra/signlink-mobile/back-end/utils/api-error"
	"github.com/devanadindra/signlink-mobile/back-end/utils/constants"
)

func GinWithCtx(ginCtx *gin.Context, ctx context.Context) *gin.Context {
	ginCtx.Request = ginCtx.
		Request.
		WithContext(
			newCombinerCtx(
				newStopperCtx(
					ginCtx.Request.Context(),
				),
				newStopperCtx(
					ctx,
				),
			),
		)
	return ginCtx

}

func GetTokenClaims(ctx context.Context) (constants.Token, error) {
	tokenVal := ctx.Value(tokenKey)
	token, ok := tokenVal.(constants.Token)
	if !ok {
		return constants.Token{}, apierror.NewWarn(http.StatusInternalServerError, "Can't get token claims")
	}
	return token, nil
}

func SetTokenClaims(ctx context.Context, token constants.Token) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func GetRequestId(ctx context.Context) *uuid.UUID {
	reqId := ctx.Value(requestIdKey)
	reqIdUUID, ok := reqId.(uuid.UUID)
	if !ok {
		return nil
	}
	return &reqIdUUID
}

func SetRequestId(ctx context.Context, requestId uuid.UUID) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}

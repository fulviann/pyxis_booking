package common

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	apierror "github.com/fulviann/pyxis_booking/back-end/utils/api-error"
	"github.com/fulviann/pyxis_booking/back-end/utils/constants"
)

func GetMetaData(ctx *gin.Context, validate *validator.Validate, allowedColumns ...string) (res *constants.FilterReq, err error) {
	limit, err := strconv.Atoi(ctx.DefaultQuery(constants.QUERY_PARAMS_LIMIT, "20"))
	if err != nil {
		return nil, apierror.NewWarn(http.StatusBadRequest, "Limit must be a number")
	}

	page, err := strconv.Atoi(ctx.DefaultQuery(constants.QUERY_PARAMS_PAGE, "1"))
	if err != nil {
		return nil, apierror.NewWarn(http.StatusBadRequest, "Page must be a number")
	}

	orderBy := ctx.DefaultQuery(constants.QUERY_PARAMS_ORDER_BY, allowedColumns[0])
	if !slices.Contains(allowedColumns, orderBy) {
		return nil, apierror.NewWarn(http.StatusBadRequest, fmt.Sprintf("Order by column '%s' is not allowed!", orderBy))
	}

	sortOrder := strings.ToLower(ctx.DefaultQuery(constants.QUERY_PARAMS_SORT_ORDER, "asc"))
	keyword := ctx.Query(constants.QUERY_PARAMS_KEYWORD)

	startCreatedAtStr := ctx.Query(constants.QUERY_PARAMS_START_CREATED_AT)
	endCreatedAtStr := ctx.Query(constants.QUERY_PARAMS_END_CREATED_AT)
	startUpdatedAtStr := ctx.Query(constants.QUERY_PARAMS_START_UPDATED_AT)
	endUpdatedAtStr := ctx.Query(constants.QUERY_PARAMS_END_UPDATED_AT)

	var startCreatedAt, endCreatedAt, startUpdatedAt, endUpdatedAt *time.Time

	if startCreatedAtStr != "" {
		tempStartCreatedAt, err := time.Parse(time.RFC3339Nano, startCreatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		startCreatedAt = &tempStartCreatedAt
	}

	if endCreatedAtStr != "" {
		tempEndCreatedAt, err := time.Parse(time.RFC3339Nano, endCreatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		endCreatedAt = &tempEndCreatedAt
	}

	if startUpdatedAtStr != "" {
		tempStartUpdatedAt, err := time.Parse(time.RFC3339Nano, startUpdatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		startUpdatedAt = &tempStartUpdatedAt
	}

	if endUpdatedAtStr != "" {
		tempEndUpdatedAt, err := time.Parse(time.RFC3339Nano, endUpdatedAtStr)
		if err != nil {
			return nil, apierror.NewWarn(http.StatusBadRequest, err.Error())
		}
		endUpdatedAt = &tempEndUpdatedAt
	}

	res = &constants.FilterReq{
		Limit:          int64(limit),
		Page:           int64(page),
		OrderBy:        orderBy,
		Keyword:        keyword,
		SortOrder:      sortOrder,
		StartCreatedAt: startCreatedAt,
		EndCreatedAt:   endCreatedAt,
		StartUpdatedAt: startUpdatedAt,
		EndUpdatedAt:   endUpdatedAt,
	}

	err = validate.Struct(res)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func Ternary[T any](condition bool, ifTrue T, ifFalse T) T {
	if condition {
		return ifTrue
	}
	return ifFalse
}

func GetValueFromPointer[V any](v *V) (res V) {
	if v != nil {
		return *v
	}
	return res
}

func ValueToPointer[V any](v V) *V {
	return &v
}

func ChunkSlice[T comparable](slice []T, chunkSize int) [][]T {
	var chunks [][]T
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// get array uniquely
func UniqueArray[T comparable](arr []T) []T {
	// Create a map to track elements that have been seen
	seen := make(map[T]bool)
	uniqueArr := make([]T, 0)

	// Iterate through the array
	for _, elem := range arr {
		// If the element is not in the map, add it to the map and the unique array
		if !seen[elem] {
			seen[elem] = true
			uniqueArr = append(uniqueArr, elem)
		}
	}

	return uniqueArr
}

func ToArrayAny(a []string) []any {
	res := make([]any, len(a))
	for i, v := range a {
		res[i] = v
	}
	return res
}

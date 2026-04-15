package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/devanadindra/signlink-mobile/back-end/utils/constants"
	"github.com/devanadindra/signlink-mobile/back-end/utils/logger"
	"github.com/gin-gonic/gin"
)

func getRequestPayload(ginCtx *gin.Context) *constants.RequestPayload {

	// get req body
	body, err := io.ReadAll(ginCtx.Request.Body)
	if err != nil {
		logger.Trace(ginCtx, "Error reading request body: %v", err)
	} else {
		// Must replace the body with original value so it can be read again
		ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	bodyStr := toValidJson(body)

	// copy header so it doest affect original header
	header := copyHeaders(ginCtx.Request.Header)
	for _, ignoredHeader := range IGNORED_HEADERS {

		if _, isMapContainsKey := header[ignoredHeader]; !isMapContainsKey {
			continue
		}

		dummyValue := make([]string, 0)
		for range header[ignoredHeader] {
			dummyValue = append(dummyValue, "***********")
		}
		header[ignoredHeader] = dummyValue

	}

	return &constants.RequestPayload{
		Body:        bodyStr,
		QueryParams: ginCtx.Request.URL.Query(),
		Headers:     header,
	}
}

func toValidJson(data []byte) (res map[string]any) {
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil
	}

	return res
}

func copyHeaders(original http.Header) map[string][]string {
	copied := make(map[string][]string)
	for key, values := range original {
		// Create a new slice for the values to ensure a deep copy
		copiedValues := make([]string, len(values))
		copy(copiedValues, values)
		copied[key] = copiedValues
	}
	return copied
}

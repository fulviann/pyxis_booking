package logger

import (
	"time"

	"github.com/fulviann/pyxis_booking/back-end/utils/constants"
)

const PACKAGE_NAME = "github.com/fulviann/pyxis_booking/back-end"

type LogPayload struct {
	Method         string
	Path           string
	StatusCode     int
	Took           time.Duration
	RequestPayload *constants.RequestPayload
}

func Setdata(env, ver string) {
	environment = env
	version = ver
}

var (
	environment = "unknown"
	version     = "unknown"
)

package middlewares

import "github.com/fulviann/pyxis_booking/back-end/utils/constants"

var IGNORED_HEADERS = []string{
	constants.AUTHORIZATION,
	constants.AUTH,
}

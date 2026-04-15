package middlewares

import "github.com/devanadindra/signlink-mobile/back-end/utils/constants"

var IGNORED_HEADERS = []string{
	constants.AUTHORIZATION,
	constants.AUTH,
}

package constant

type ErrorEnums struct {
	ErrorCode    int
	ErrorMessage string
}

var LogPattern = "[%s] [%s] - %s"

var ErrorConstant = map[string]ErrorEnums{
	"SUCCESS": {
		ErrorCode:    0,
		ErrorMessage: "Success.",
	},
	"INTERNAL_FAILURE": {
		ErrorCode:    -1,
		ErrorMessage: "An error has been occurred, please try again later.",
	},
	"QUERY_ERROR": {
		ErrorCode:    1,
		ErrorMessage: "Query error.",
	},
	"CREATE_DUPLICATE_USER": {
		ErrorCode:    2,
		ErrorMessage: "User already exist.",
	},
	"JSON_BINDING_ERROR": {
		ErrorCode:    3,
		ErrorMessage: "Json binding error.",
	},
	"AUTHENTICATE_FAILURE": {
		ErrorCode:    4,
		ErrorMessage: "Authenticate fail.",
	},
	"UNAUTHORIZED": {
		ErrorCode:    5,
		ErrorMessage: "Unauthorized.",
	},
}

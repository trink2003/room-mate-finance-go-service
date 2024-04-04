package constant

type ErrorEnums struct {
	ErrorCode    int
	ErrorMessage string
}

const LogPattern = "[%s] [%s] - %s"

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
	"DATA_FORMAT_ERROR": {
		ErrorCode:    6,
		ErrorMessage: "Data format error.",
	},
	"USER_NOT_EXISTED": {
		ErrorCode:    7,
		ErrorMessage: "User not existed.",
	},
	"INVALID_NUMBER_OF_USER": {
		ErrorCode:    8,
		ErrorMessage: "The number of users in the same room must be greater than 2.",
	},
	"INVALID_USER_TO_PAID_LIST": {
		ErrorCode:    9,
		ErrorMessage: "The buyer must not be on the list of payers.",
	},
}

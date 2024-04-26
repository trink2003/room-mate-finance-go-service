package constant

type LogLevelType string

const (
	Info  LogLevelType = "INFO"
	Warn  LogLevelType = "WARN"
	Error LogLevelType = "ERROR"
	Debug LogLevelType = "DEBUG"
)

type ErrorEnums struct {
	ErrorCode    int
	ErrorMessage string
}

const BaseApiPath = "/roommate/api/v1"

const LogFileFolder = "./service_log/"
const LogFileLocation = "room_mate_finance_service_log_%d_%d_%d.log"
const DeltaPositive = 0.5
const DeltaNegative = -0.5
const YyyyMmDdHhMmSsFormat = "2006-01-02 15:04:05"
const AscKeyword = "ASC"
const DescKeyword = "DESC"
const EmptyString = ""
const (
	ContentTypeBinary    = "application/octet-stream"
	ContentTypeForm      = "application/x-www-form-urlencoded"
	ContentTypeJSON      = "application/json"
	ContentTypeHTML      = "text/html; charset=utf-8"
	ContentTypeText      = "text/plain; charset=utf-8"
	ContentTypeIconImage = "image/x-icon"
)

var SensitiveField = [...]string{"password", "jwt", "token", "client_secret", "Authorization", "x-api-key"} // [...] instead of []: it ensures you get a (fixed size) array instead of a slice. So the values aren't fixed but the size is.
var ValidMethod = []string{"GET", "POST", "PUT", "DELETE"}

type LogKey string

const UsernameLogKey LogKey = "username"
const TraceIdLogKey LogKey = "traceId"
const LogPattern = "[%s] [%s] 👉️ \t%s"

var (
	Success = ErrorEnums{
		ErrorCode:    0,
		ErrorMessage: "Success",
	}
	InternalFailure = ErrorEnums{
		ErrorCode:    -1,
		ErrorMessage: "An error has been occurred, please try again later",
	}
	PageNotFound = ErrorEnums{
		ErrorCode:    -2,
		ErrorMessage: "You're consuming an unknow endpoint, please check your url (404 Page Not Found)",
	}
	MethodNotAllowed = ErrorEnums{
		ErrorCode:    -3,
		ErrorMessage: "This url is configured method that not match with your current method, please check again (405 Method Not Allowed)",
	}
	QueryError = ErrorEnums{
		ErrorCode:    1,
		ErrorMessage: "Query error",
	}
	CreateDuplicateUser = ErrorEnums{
		ErrorCode:    2,
		ErrorMessage: "User already exist",
	}
	JsonBindingError = ErrorEnums{
		ErrorCode:    3,
		ErrorMessage: "Json binding error",
	}
	AuthenticateFailure = ErrorEnums{
		ErrorCode:    4,
		ErrorMessage: "Authenticate fail",
	}
	Unauthorized = ErrorEnums{
		ErrorCode:    5,
		ErrorMessage: "Unauthorized",
	}
	DataFormatError = ErrorEnums{
		ErrorCode:    6,
		ErrorMessage: "Data format error",
	}
	UserNotExisted = ErrorEnums{
		ErrorCode:    7,
		ErrorMessage: "User not existed",
	}
	InvalidNumberOfUser = ErrorEnums{
		ErrorCode:    8,
		ErrorMessage: "The number of users in the same room must be greater than 2",
	}
	InvalidUserToPaidList = ErrorEnums{
		ErrorCode:    9,
		ErrorMessage: "The buyer must not be on the list of payers",
	}
	ExpenseDeleteNotSuccess = ErrorEnums{
		ErrorCode:    10,
		ErrorMessage: "An error occurred while deleting daily spending data",
	}
	ExpenseActiveNotSuccess = ErrorEnums{
		ErrorCode:    11,
		ErrorMessage: "An error occurred while activating daily spending data",
	}
	Forbidden = ErrorEnums{
		ErrorCode:    12,
		ErrorMessage: "You don't have permission to perform this action",
	}
	RoomHasBeenExisted = ErrorEnums{
		ErrorCode:    13,
		ErrorMessage: "Room has been existed",
	}
	RoomDoesNotExist = ErrorEnums{
		ErrorCode:    14,
		ErrorMessage: "Room does not exist",
	}
	RoomStillHavePeople = ErrorEnums{
		ErrorCode:    15,
		ErrorMessage: "The room you want to remove still have people live in there",
	}
	RoleNotExist = ErrorEnums{
		ErrorCode:    16,
		ErrorMessage: "Can not find role for user",
	}
	DeleteDefaultRoomError = ErrorEnums{
		ErrorCode:    17,
		ErrorMessage: "Can not delete default room",
	}
	ActionCannotPerformOnYourself = ErrorEnums{
		ErrorCode:    18,
		ErrorMessage: "You cannot perform this action to your self",
	}
	EmptyRoomError = ErrorEnums{
		ErrorCode:    19,
		ErrorMessage: "Room currenty does not have any memeber",
	}
	ChangeRoomForAllMemberError = ErrorEnums{
		ErrorCode:    20,
		ErrorMessage: "Change room for all member error",
	}
)

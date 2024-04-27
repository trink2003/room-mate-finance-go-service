package log

import (
	"context"
	"fmt"
	"os"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/utils/splunk/v2"
	"time"

	"github.com/charmbracelet/log"
)

func WithLevel(level constant.LogLevelType, ctx context.Context, content string) {

	// ensure that ctx is never nil
	if ctx == nil {
		ctx = context.Background()
		ctx = context.WithValue(ctx, constant.UsernameLogKey, "nil ctx input")
		ctx = context.WithValue(ctx, constant.TraceIdLogKey, "nil ctx input")
	}

	timeZoneLocation, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		return
	}
	currentTimestamp := time.Now().In(timeZoneLocation)
	usernameFromContext := ctx.Value(constant.UsernameLogKey)
	traceIdFromContext := ctx.Value(constant.TraceIdLogKey)
	username := constant.EmptyString
	traceId := constant.EmptyString
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	// fmt.Println(strings.Compare(string(level), string(constant.LogLevelType("INFO"))))
	var message = fmt.Sprintf(
		constant.LogPattern,
		traceId,
		username,
		content,
	)
	switch level {
	case constant.Info:
		log.Info(
			message,
		)
	case constant.Warn:
		log.Warn(
			message,
		)
	case constant.Error:
		log.Error(
			message,
		)
	case constant.Debug:
		log.Debug(
			message,
		)
	default:
		log.Info(
			message,
		)
	}

	host, token, source, sourcetype, index, splunkInfoIsFullSetInEnv := GetSplunkInformationFromEnvironment()

	if splunkInfoIsFullSetInEnv {
		splunkClient := splunk.NewClient(
			nil,
			host,
			token,
			source,
			sourcetype,
			index,
		)
		err := splunkClient.Log(
			message,
		)
		if err != nil {
			log.Error(err)
		}
	}

	appendLogToFileError := AppendLogToFile(
		fmt.Sprintf(
			"%s: %s - %s\n",
			currentTimestamp.Format(constant.YyyyMmDdHhMmSsFormat),
			string(level),
			message,
		),
	)
	if appendLogToFileError != nil {
		log.Error(fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			fmt.Sprintf("An error has been occurred when appending log to file: %s", appendLogToFileError.Error()),
		))
	}
}

// GetSplunkInformationFromEnvironment
// SPLUNK_HOST: "https://{your-splunk-URL}:8088/services/collector",
// SPLUNK_TOKEN: "{your-token}",
// SPLUNK_SOURCE: "{your-source}",
// SPLUNK_SOURCETYPE: "{your-sourcetype}",
// SPLUNK_INDEX: "{your-index}",
func GetSplunkInformationFromEnvironment() (host string, token string, source string, sourcetype string, index string, splunkInfoIsFullSetInEnv bool) {
	var splunkHost, isSplunkHostSet = os.LookupEnv("SPLUNK_HOST")
	var splunkToken, isSplunkTokenSet = os.LookupEnv("SPLUNK_TOKEN")
	var splunkSource, isSplunkSourceSet = os.LookupEnv("SPLUNK_SOURCE")
	var splunkSourcetype, isSplunkSourcetypeSet = os.LookupEnv("SPLUNK_SOURCETYPE")
	var splunkIndex, isSplunkIndexSet = os.LookupEnv("SPLUNK_INDEX")
	if !isSplunkHostSet && !isSplunkTokenSet && !isSplunkSourceSet && !isSplunkSourcetypeSet && !isSplunkIndexSet {
		return "", "", "", "", "", false
	}
	return splunkHost, splunkToken, splunkSource, splunkSourcetype, splunkIndex, true
}

func AppendLogToFile(log string) error {
	timeZoneLocation, timeLoadLocationErr := time.LoadLocation("Asia/Ho_Chi_Minh")
	if timeLoadLocationErr != nil {
		return timeLoadLocationErr
	}
	currentTimestamp := time.Now().In(timeZoneLocation)

	logFileName := fmt.Sprintf(constant.LogFileLocation, currentTimestamp.Year(), int(currentTimestamp.Month()), currentTimestamp.Day())

	// check if log folder is existed or not
	if _, directoryStatusError := os.Stat(constant.LogFileFolder); os.IsNotExist(directoryStatusError) {
		makeDirectoryAllError := os.MkdirAll(constant.LogFileFolder, 0755)
		if makeDirectoryAllError != nil {
			return makeDirectoryAllError
		}
	}

	// +-----+---+--------------------------+
	// | rwx | 7 | Read, write and execute  |
	// | rw- | 6 | Read, write              |
	// | r-x | 5 | Read, and execute        |
	// | r-- | 4 | Read,                    |
	// | -wx | 3 | Write and execute        |
	// | -w- | 2 | Write                    |
	// | --x | 1 | Execute                  |
	// | --- | 0 | no permissions           |
	// +------------------------------------+

	// +------------+------+-------+
	// | Permission | Octal| Field |
	// +------------+------+-------+
	// | rwx------  | 0700 | User  |
	// | ---rwx---  | 0070 | Group |
	// | ------rwx  | 0007 | Other |
	// +------------+------+-------+
	// O_RDONLY: It opens the file read-only.
	// O_WRONLY: It opens the file write-only.
	// O_RDWR: It opens the file read-write.
	// O_APPEND: It appends data to the file when writing.
	// O_CREATE: It creates a new file if none exists.
	file, openFileError := os.OpenFile(constant.LogFileFolder+logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if openFileError != nil {
		return openFileError
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("An error has been occurred when defer closing file: %v", err)
		}
	}(file)

	_, writeStringToFileError := file.WriteString(log)

	if writeStringToFileError != nil {
		return writeStringToFileError
	}
	return nil
}

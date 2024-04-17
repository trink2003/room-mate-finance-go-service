package log

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"room-mate-finance-go-service/constant"
	"room-mate-finance-go-service/utils/splunk/v2"
	"strings"
)

func WithLevel(level constant.LogLevelType, ctx context.Context, content string) {
	usernameFromContext := ctx.Value("username")
	traceIdFromContext := ctx.Value("traceId")
	username := ""
	traceId := ""
	if usernameFromContext != nil {
		username = usernameFromContext.(string)
	}
	if traceIdFromContext != nil {
		traceId = traceIdFromContext.(string)
	}
	fmt.Println(strings.Compare(string(level), string(constant.LogLevelType("INFO"))))
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
		break
	case constant.Warn:
		log.Warn(
			message,
		)
		break
	case constant.Error:
		log.Error(
			message,
		)
		break
	default:
		log.Info(
			message,
		)
		break
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
	if isSplunkHostSet == false && isSplunkTokenSet == false && isSplunkSourceSet == false && isSplunkSourcetypeSet == false && isSplunkIndexSet == false {
		return "", "", "", "", "", false
	}
	return splunkHost, splunkToken, splunkSource, splunkSourcetype, splunkIndex, true
}

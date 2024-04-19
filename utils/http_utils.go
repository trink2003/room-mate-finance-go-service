package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"room-mate-finance-go-service/constant"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
)

func ConsumeApi(
	ctx context.Context,
	url string,
	method string,
	header map[string]string,
	payload string,
	isVerifySsl bool,
) (string, error) {
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

	if slices.Contains(constant.ValidMethod, method) == false {
		return "", errors.New("invalid method")
	}
	var client *http.Client
	if isVerifySsl {
		client = &http.Client{}
	} else {
		customTransport := http.DefaultTransport.(*http.Transport).Clone()
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		client = &http.Client{Transport: customTransport}
	}
	req, err := http.NewRequest(method, url, strings.NewReader(payload))
	if err != nil {
		log.Error(
			fmt.Sprintf(
				constant.LogPattern,
				traceId,
				username,
				"ConsumeApi - http.NewRequest - error: "+err.Error(),
			),
		)
		return "", err
	}

	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			curlBuilder(url, payload, header),
		),
	)

	for k, v := range header {
		req.Header.Add(k, v)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Error(
			fmt.Sprintf(
				constant.LogPattern,
				traceId,
				username,
				"ConsumeApi - client.Do - error: "+err.Error(),
			),
		)
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(res.Body)

	resHeader := map[string][]string(res.Header)

	headerString := ""

	for k, v := range resHeader {
		if IsSensitiveField(k) {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, "***")
		} else {
			headerString += fmt.Sprintf("\n\t\t- %s: %s", k, strings.Join(v, ", "))
		}
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error(
			fmt.Sprintf(
				constant.LogPattern,
				traceId,
				username,
				"ConsumeApi - io.ReadAll - error: "+err.Error(),
			),
		)
		return "", err
	}
	result := string(body)
	log.Info(
		fmt.Sprintf(
			constant.LogPattern,
			traceId,
			username,
			fmt.Sprintf(
				"\t- status: %s\n\t- header: %s\n\t- payload: %s",
				res.Status,
				headerString,
				result,
			),
		),
	)

	return result, nil

}

func curlBuilder(url string, payload string, header map[string]string) string {
	curlCommand := "curl "
	curlCommand += "'" + url + "' "
	for k, v := range header {
		curlCommand += "-H '" + k + ": " + v + "' "
	}
	if payload != "" {
		curlCommand += "-X POST -d '" + payload + "'"
	} else {
		curlCommand += "-X GET"
	}
	return curlCommand
}

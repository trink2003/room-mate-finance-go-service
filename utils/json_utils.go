package utils

import (
	"encoding/json"
	"room-mate-finance-go-service/constant"
	"strings"
)

func StructToJson(anyStruct any) string {
	result, err := json.Marshal(anyStruct)
	if err != nil {
		return ""
	}
	return string(result)
}

func JsonToStruct[T any](jsonString string, anyStruct *T) {
	err := json.Unmarshal([]byte(jsonString), anyStruct)
	if err != nil {
		return
	}
}

func ByteJsonToStruct[T any](jsonString []byte, anyStruct *T) {
	err := json.Unmarshal(jsonString, anyStruct)
	if err != nil {
		return
	}
}

func SortMapToString(inputMap map[string]string) string {
	result := ""

	if inputMap == nil || len(inputMap) == 0 {
		return result
	}

	for k, v := range inputMap {
		sort := ""
		if v != constant.AscKeyword && v != constant.DescKeyword {
			sort = constant.DescKeyword
		} else {
			sort = v
		}
		result += k + " " + sort + ", "
	}
	return strings.TrimSuffix(result, ", ")
}

func HideSensitiveJsonField(inputJson string) string {
	element := strings.Split(inputJson, "\"")
	for i := range element {
		currentField := element[i]
		var colon string
		if (len(element) == 0) || (i+1 > len(element)-1) {
			continue
		}
		colon = element[i+1]
		if IsSensitiveField(currentField) {
			if strings.Contains(strings.Trim(colon, " "), ":") {
				element[i+2] = "***"
			}
		} else if i+2 < len(element) && len(element[i+2]) > 1000 {
			element[i+2] = element[i+2][:50] + "..." + element[i+2][len(element[i+2])-50:]
		}
	}
	return strings.Join(element, "\"")
}

func IsSensitiveField(input string) bool {
	for _, e := range constant.SensitiveField {
		if strings.Contains(strings.ToLower(e), strings.ToLower(input)) || strings.Contains(strings.ToLower(input), strings.ToLower(e)) {
			return true
		}
	}
	return false
}

package viewer

import (
	"strings"
)

func Parse(htmlTemplate string, metadata map[string]string) string {
	for key, value := range metadata {
		replace := "{{" + key + "}}"
		htmlTemplate = strings.Replace(htmlTemplate, replace, value, -1)
	}

	return htmlTemplate
}

package utils

import "strings"

func ParseTagSetting(str string, sep string) map[string]string {
	settings := map[string]string{}
	names := strings.Split(str, sep)

	for i := 0; i < len(names); i++ {
		j := i

		values := strings.Split(names[j], ":")
		k := strings.TrimSpace(values[0])

		if len(values) >= 2 {
			settings[k] = strings.Join(values[1:], ":")
		} else if k != "" {
			settings[k] = ""
		}
	}

	return settings
}

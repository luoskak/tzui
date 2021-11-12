package utils

import "strings"

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms            = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer    *strings.Replacer
	commonInitialisCasedReplacer *strings.Replacer
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, strings.Title(strings.ToLower(initialism)))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
	commonInitialismsForCasedReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForCasedReplacer = append(commonInitialismsForCasedReplacer, initialism, casedCommonInitialisms(initialism))
	}
	commonInitialisCasedReplacer = strings.NewReplacer(commonInitialismsForCasedReplacer...)
}

func casedCommonInitialisms(initialism string) string {
	lower := strings.ToLower(initialism)
	old := []rune(lower)
	top := old[0]
	top += 'A' - 'a'
	return string(append([]rune{top}, old[1:]...))
}

func CasedName(name string) string {
	name = commonInitialisCasedReplacer.Replace(name)
	old := []rune(name)
	top := old[0]
	top += 'a' - 'A'
	return string(append([]rune{top}, old[1:]...))
}

func SnackedName(name string) string {
	if name == "" {
		return ""
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase

	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}
	ret := buf.String()
	return ret
}

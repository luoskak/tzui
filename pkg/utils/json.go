package utils

import (
	jsoniter "github.com/json-iterator/go"
)

type JSONExtension struct {
	jsoniter.DummyExtension
}

func (extension *JSONExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		var ns []string
		for _, oldName := range field.ToNames {
			oldName = commonInitialisCasedReplacer.Replace(oldName)
			old := []rune(oldName)
			top := old[0]
			if top <= 'z' && top >= 'a' {
				ns = append(ns, oldName)
				continue
			}
			top += 'a' - 'A'
			ns = append(ns, string(append([]rune{top}, old[1:]...)))
		}
		field.FromNames = ns
		field.ToNames = ns
	}
}

// JsonAPI ignore json naming tag json:"ignore,omitempty"
// UserName => userName
// ID => id
// UserID => userId
var JsonAPI jsoniter.API

func init() {
	JsonAPI = jsoniter.Config{
		EscapeHTML:    true,
		CaseSensitive: true,
	}.Froze()
	JsonAPI.RegisterExtension(&JSONExtension{})
}

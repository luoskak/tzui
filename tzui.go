package tzui

import (
	"errors"
	"strings"
)

var ErrUnsupportedModelType = errors.New("unsupported model type")

func NewTzuiPageBuilder(root, sub string, name string) *pageBuilder {
	sub = strings.TrimSpace(sub)
	if root != "" {
		sub = strings.TrimSpace(root) + "/" + sub
	}
	page := &pageBuilder{
		name: name,
		sub:  sub,
	}
	return page
}

func NewTzuiControllerBuilder(root, sub string, name, desc string) *controllerBuilder {
	root = "/" + strings.TrimPrefix(root, "/")
	c := &controllerBuilder{
		root: root,
		sub:  strings.TrimSpace(sub),
		name: name,
		desc: desc,
	}
	c.tagManager = &TagManager{
		cbuild: c,
		tags:   make(map[string]Tag),
	}
	return c
}

func NewTzuiDictionaryBuilder(root, sub string) *dictionaryBuilder {
	root = "/" + strings.TrimPrefix(root, "/")
	return &dictionaryBuilder{
		root: root,
		sub:  sub,
		dics: make(map[string]*dictionary),
	}
}

package tzui

import (
	"context"
	"fmt"
	"strings"
)

type controllerBuilder struct {
	root, sub  string
	name, desc string
	builders   []*pageBuilder
	tagManager *TagManager
	Error      error
}

func (cbuild *controllerBuilder) AddPage(name string, sub string) *pageBuilder {
	sub = strings.TrimSpace(sub)
	if sub == "" {
		panic("page sub can not be empty")
	}
	page := &pageBuilder{
		name:   name,
		sub:    sub,
		cbuild: cbuild,
	}
	cbuild.builders = append(cbuild.builders, page)
	return page
}

func (cbuild controllerBuilder) ResolveMethods(register func(route, name, desc string, binder func(body []byte) (interface{}, error), handle func(ctx context.Context, req interface{}) (interface{}, error))) {
	for _, pbuild := range cbuild.builders {
		var pageURL string
		if cbuild.sub != "" {
			pageURL = strings.Join([]string{
				cbuild.sub,
				"page",
				pbuild.sub,
			}, "/")
		} else {
			pageURL = "page" + "/" + pbuild.sub
		}

		register("/"+pageURL, cbuild.name, pbuild.name, nil, pbuild.Handle)

		for _, c := range pbuild.components {
			var rootPath string
			if cbuild.sub != "" {
				rootPath = strings.Join(
					[]string{
						cbuild.sub,
						pbuild.sub,
						c.cmp.ComponentName(),
					},
					"/",
				)
			} else {
				rootPath = pbuild.sub + "/" + c.cmp.ComponentName()
			}

			for _, source := range c.sources {
				if !source.initialized {
					url := cbuild.root + "/" + rootPath + "/" + source.URL
					source.URL = url
				}
				source.initialized = true
				register(strings.TrimPrefix(source.URL, cbuild.root), pbuild.name, source.Name, source.Binder, source.Handler)
			}
		}
	}
}

func (cbuild *controllerBuilder) BindDictionary(dbuild *dictionaryBuilder) {
	tag := DictionaryTag{
		modifiedTags: make(map[string]bool),
		tm:           cbuild.tagManager,
	}
	for name := range dbuild.dics {
		tag.modifiedTags[name] = true
	}
	cbuild.tagManager.AddTag("dictionary", tag)
}

func (cbuild *controllerBuilder) addError(err error) error {
	if cbuild.Error == nil {
		cbuild.Error = err
	} else if err != nil {
		cbuild.Error = fmt.Errorf("%v; %w", cbuild.Error, err)
	}
	return cbuild.Error
}

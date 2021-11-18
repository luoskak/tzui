package tzui

import (
	"context"
	"fmt"
	"strings"
)

// TODO: support for pwd
type dictionaryBuilder struct {
	root, sub string
	dics      map[string]*dictionary
}

type dictionary struct {
	name    string
	handler DictionaryHandler
}

type DictionaryHandler func(ctx context.Context, req *TzDictionaryRequest) (interface{}, error)

func (dbuild *dictionaryBuilder) AddDictionary(name string, handler DictionaryHandler) *dictionaryBuilder {
	_, duplicated := dbuild.dics[name]
	if duplicated {
		panic(fmt.Sprintf("duplicated add dictionary %s for sub %s", name, dbuild.sub))
	}
	dbuild.dics[name] = &dictionary{
		name:    name,
		handler: handler,
	}
	return dbuild
}

func (dbuild *dictionaryBuilder) ResolveMethods(register func(route, name, desc string, binder func(body []byte) (interface{}, error), handle func(ctx context.Context, req interface{}) (interface{}, error))) {
	var dics []TzDictionary
	for _, d := range dbuild.dics {
		dic := TzDictionary{
			Name: d.name,
			URL: dbuild.root + strings.Join([]string{
				dbuild.sub,
				"dic",
				d.name,
			}, "/"),
		}
		dics = append(dics, dic)
	}
	var routePath string
	if dbuild.sub != "" {
		routePath = "/" + dbuild.sub
	}
	register(routePath+"/dic", "字典服务", "字典服务", nil, func(ctx context.Context, req interface{}) (interface{}, error) {
		return dics, nil
	})
	for _, dic := range dics {
		d := dbuild.dics[dic.Name]
		register(strings.TrimPrefix(dic.URL, dbuild.root), d.name, "字典", (&TzDictionaryRequest{}).Bind, func(ctx context.Context, req interface{}) (interface{}, error) {
			return d.handler(ctx, req.(*TzDictionaryRequest))
		})
	}
}

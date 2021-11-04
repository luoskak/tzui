package tzui

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"reflect"
	"strings"
	"sync"

	"gitlab.com/tz/tzui/pkg/utils"
)

var ErrUnsupportedModelType = errors.New("unsupported model type")

var jsonAPI = utils.JsonAPI

var (
	cacheComponentStore = new(sync.Map)
)

type ControllerBuilder struct {
	root     string
	builders []*pageBuilder
}

func NewTzuiControllerBuilder(root string) *ControllerBuilder {
	return &ControllerBuilder{
		root: root,
	}
}

type pageBuilder struct {
	name       string
	components []ITzComponent
	sourceURLs []string
	handlers   []Handler
}

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type ITzComponentBuilder interface {
	Model() interface{}
	Build() ITzComponent
	ComponentName() string
}

type HasSourceComponentBuilder interface {
	SourceURL() string
	Handler(ctx context.Context, req interface{}) (interface{}, error)
}

func (pbuild *pageBuilder) build() *TzPage {
	return &TzPage{
		TzComponent: TzComponent{
			Name:     pbuild.name,
			Children: pbuild.components,
		},
	}
}

func (pbuild *pageBuilder) append(source string, cmp ITzComponent, handle Handler) {
	pbuild.sourceURLs = append(pbuild.sourceURLs, "")
	pbuild.components = append(pbuild.components, cmp)
	pbuild.handlers = append(pbuild.handlers, handle)
}

func (pbuild *pageBuilder) AddTzComponent(cbuild ITzComponentBuilder) *pageBuilder {
	var (
		model     = cbuild.Model()
		sourceURL = ""
		handle    Handler
	)

	if model == nil {
		// TODO no model component example
		cmp := cbuild.Build()
		// sourceURL为空
		pbuild.append(sourceURL, cmp, handle)
		return pbuild
	}
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	if modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			panic(fmt.Errorf("%w: %+v", ErrUnsupportedModelType, model))
		}
		panic(fmt.Errorf("%w: %s.%s", ErrUnsupportedModelType, modelType.PkgPath(), modelType.Name()))
	}

	if app, ok := cbuild.(HasSourceComponentBuilder); ok {
		sourceURL = app.SourceURL()
		handle = app.Handler
	}
	pageName := fmt.Sprintf("%s.%s.%s", modelType.PkgPath(), modelType.Name(), cbuild.ComponentName())
	if v, cached := cacheComponentStore.Load(pageName); cached {
		pbuild.append(sourceURL, v.(ITzComponent), handle)
		return pbuild
	}

	var fields []*Field
	for i := 0; i < modelType.NumField(); i++ {
		if fieldStruct := modelType.Field(i); ast.IsExported(fieldStruct.Name) {
			field := Field{
				fieldName: fieldStruct.Name,
				Tags:      utils.ParseTagSetting(fieldStruct.Tag.Get("tzui"), ";"),
			}
			fields = append(fields, &field)
		}
	}
	cmp := cbuild.Build()
	if app, ok := cmp.(IParseTag); ok {
		app.ParseTag(fields)
	}

	cacheComponentStore.Store(pageName, cmp)
	pbuild.append(sourceURL, cmp, handle)
	return pbuild
}

type TzPage struct {
	TzComponent
}

func (p *TzPage) Binder(body []byte) (interface{}, error) {
	err := jsonAPI.Unmarshal(body, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (cbuild *ControllerBuilder) AddPage(name string) *pageBuilder {
	page := &pageBuilder{
		name: name,
	}
	cbuild.builders = append(cbuild.builders, page)
	return page
}

func (cbuild ControllerBuilder) ResolveMethods(register func(route, name, desc string, binder func(body []byte) (interface{}, error), handle func(ctx context.Context, req interface{}) (interface{}, error))) {
	for _, pbuild := range cbuild.builders {
		pageURL := strings.Join([]string{
			cbuild.root,
			"page",
			pbuild.name,
		}, "/")
		page := pbuild.build()
		register(pageURL, pbuild.name, "动态路由配置", page.Binder, pbuild.handle)
		for i, c := range pbuild.components {

			if app, ok := c.(HasSourceComponent); ok {
				ps := []string{
					cbuild.root,
					pbuild.name,
					c.ComponentName(),
				}
				if pbuild.sourceURLs[i] != "" {
					ps = append(ps, pbuild.sourceURLs[i])
				}
				url := strings.Join(ps, "/")
				app.SetSource(url)
				// 获取模型的路径
				register(url, pbuild.sourceURLs[i], "模型", func(body []byte) (interface{}, error) {
					req := c.ReqStr()
					err := jsonAPI.Unmarshal(body, req)
					if err != nil {
						return nil, err
					}
					return req, nil
				}, pbuild.handlers[i])
			}

		}
	}
}

func (pbuild *pageBuilder) handle(ctx context.Context, req interface{}) (interface{}, error) {
	var page *TzPage
	if req == nil {
		page = pbuild.build()
	} else {
		page = req.(*TzPage)
		page.Children = pbuild.components
	}
	return page, nil
}

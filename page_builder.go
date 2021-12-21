package tzui

import (
	"context"
	"fmt"
	"go/ast"
	"reflect"
	"sync"

	"gitlab.com/tz/tzui/pkg/utils"
)

var (
	cacheComponentStore = new(sync.Map)
)

type pageBuilder struct {
	// sub path
	cbuild     *controllerBuilder
	sub        string
	name       string
	components []*pageBuilderCom
}

type pageBuilderCom struct {
	cmp     ITzComponent
	sources []*TzSource
}

func (pbuild *pageBuilder) build() *TzPage {
	var children []ITzComponent
	for _, c := range pbuild.components {
		children = append(children, c.cmp)
	}
	return &TzPage{
		TzComponent: TzComponent{
			Name:     pbuild.name,
			Children: children,
		},
	}
}

// 添加并执行组件创建器
func (pbuild *pageBuilder) AddTzComponent(cbuild ITzComponentBuilder) *pageBuilder {
	var (
		model = cbuild.Model()
	)

	if model == nil {
		// TODO no model component example
		cmp := cbuild.Build()
		// sourceURL为空
		pbuild.components = append(pbuild.components, &pageBuilderCom{cmp: cmp, sources: cbuild.Sources()})
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

	pageName := fmt.Sprintf("%s.%s.%s", modelType.PkgPath(), modelType.Name(), cbuild.ComponentName())
	if v, cached := cacheComponentStore.Load(pageName); cached {
		pbuild.components = append(pbuild.components, &pageBuilderCom{cmp: v.(ITzComponent), sources: cbuild.Sources()})
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
		var tm *TagManager
		if pbuild.cbuild != nil {
			tm = pbuild.cbuild.tagManager
		}
		app.ParseTag(fields, tm)
	}

	cacheComponentStore.Store(pageName, cmp)
	pbuild.components = append(pbuild.components, &pageBuilderCom{cmp: cmp, sources: cbuild.Sources()})
	return pbuild
}

type TzPage struct {
	TzComponent
}

func (pbuild *pageBuilder) Handle(ctx context.Context, req interface{}) (interface{}, error) {
	return pbuild.build(), nil
}

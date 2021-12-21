package devui

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"gitlab.com/tz/tzui"
	"gitlab.com/tz/tzui/pkg/utils"
)

var json = utils.JsonAPI

var (
	_ tzui.IParseTag           = &DataTableTzComponent{}
	_ tzui.ITzComponentBuilder = &dataTableTzComponentBuilder{}
)

// tzui tags:
// - header 表头名
// - fieldType 单元格格式 text
// - fieldLeft 左悬浮长度，用于固定列 0px 则第一列固定
// - width 单元格固定长度
type DataTableTzComponent struct {
	tzui.TzComponent
	Scrollable             bool   `json:"scrollable"`
	Type                   string `json:"type"`
	TableHeight            string `json:"tableHeight"`
	VirtualScroll          bool   `json:"virtualScroll"`
	FixHeader              bool   `json:"fixHeader"`
	ContainFixHeaderHeight bool   `json:"containFixHeaderHeight"`
	AutoAddTotalRow        bool
	TableWidthConfig       []TableWidthConfig `json:"tableWidthConfig"`
	TableOptions           []TableOption      `json:"tableOptions"`
}

type TableWidthConfig struct {
	Field    string `json:"field"`
	Width    string `json:"width"`
	widthNum int
}

type TableOption struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	// typ 'text' | 'enum' 'dic' enum不显示value,dic显示[id]
	FieldType   string `json:"fieldType"`
	EnumName    string
	Placeholder string //当数据为空时填充值
	// form search or datePicker, dateRangePicker
	Form string
	// exp 120px
	FixedLeft     string `json:"fixedLeft"`
	ResizeEnabled bool   `json:"resizeEnable"`
	Sortable      bool
	Sort          string
	// 百分比
	Percent string
}

type DataTableSourceRequest struct {
	Where map[string]interface{} `json:"where"`
	// sort config exp. [][]string{[]string{'name','desc'}}
	// 使用前必须在结构项加上tag sortable 否则会被前端过滤
	Sort    [][]string
	Total   int64 `json:"total"`
	PerPage int   `json:"perPage"`
	Page    int   `json:"page"`
}

func (req *DataTableSourceRequest) Data(data interface{}) *DataTableSourceResponse {
	return &DataTableSourceResponse{
		Total:   req.Total,
		PerPage: req.PerPage,
		Page:    req.Page,
		Data:    data,
	}
}

// Slice directly return when data is nil and data type is not slice
func (req *DataTableSourceRequest) Slice(data interface{}) (res *DataTableSourceResponse) {
	res = &DataTableSourceResponse{
		Total:   req.Total,
		PerPage: req.PerPage,
		Page:    req.Page,
	}
	if data == nil {
		return
	}
	if req.Page == 0 || req.PerPage == 0 {
		return
	}
	rv := reflect.ValueOf(data)
	rt := rv.Type()
	if rt.Kind() != reflect.Slice {
		panic("slice kind required")
	}
	total := rv.Len()
	from := (req.Page - 1) * req.PerPage
	to := req.Page * req.PerPage
	if from > total {
		return
	}
	if from < total && to > total {
		res.Data = rv.Slice(from, total).Interface()
		return
	}
	res.Data = rv.Slice(from, to).Interface()
	return
}

func (req *DataTableSourceRequest) WhereGet(snackedName string, equality string) (interface{}, bool) {
	if snackedName == "" {
		panic("snacked name can not be empty")
	}
	if equality == "" {
		panic("equality can not be empty")
	}
	fwi, ok := req.Where[snackedName]
	if !ok {
		return nil, false
	}
	fw := fwi.(map[string]interface{})
	equality = strings.ToUpper(equality)
	if equality == "OR" || equality == "AND" {
		panic("unsupported equality")
	}
	wvi, ok := fw[equality]
	if !ok {
		return nil, false
	}
	return wvi, true
}

func (req *DataTableSourceRequest) WhereFind(snackedNames ...string) map[string]interface{} {
	if len(snackedNames) == 0 {
		return req.Where
	}
	return findWhere(req.Where, snackedNames...)
}

func findWhere(src map[string]interface{}, snackedNames ...string) map[string]interface{} {
	where := make(map[string]interface{})

	for n, v := range src {
		used := false
		for _, fsn := range snackedNames {
			if n == fsn {
				used = true
				break
			}
		}
		if !used {
			continue
		}
		subWhere := make(map[string]interface{})
		for sn, sv := range v.(map[string]interface{}) {
			if sn == "OR" ||
				sn == "AND" {
				subWhere[sn] = findWhere(sv.(map[string]interface{}), snackedNames...)
				continue
			}
			subWhere[sn] = sv
		}
		where[n] = subWhere
	}

	return where
}

func (*DataTableSourceRequest) Bind(body []byte) (interface{}, error) {
	req := new(DataTableSourceRequest)
	err := json.Unmarshal(body, req)
	// 将sort的字段都转换成符合数据库命名规范
	var sorts [][]string
	for _, sort := range req.Sort {
		if len(sort) == 0 {
			continue
		}
		name := utils.SnackedName(sort[0])
		direction := "ASC"
		if len(sort) == 2 && sort[1] != "" {
			direction = strings.ToUpper(sort[1])
		}
		if direction != "ASC" && direction != "DESC" {
			return nil, fmt.Errorf("sort %s set but %s is not support", sort[0], sort[1])
		}
		sorts = append(sorts, []string{name, direction})
	}
	req.Sort = sorts
	// 将where的字段都转换成符合数据库命名规范
	where := convertWhere(req.Where)
	req.Where = where

	return req, err
}

func convertWhere(src map[string]interface{}) map[string]interface{} {
	where := make(map[string]interface{})
	for n, v := range src {
		sub, ok := v.(map[string]interface{})
		if !ok {
			// 不符合要求的过滤了
			continue
		}
		subWhere := make(map[string]interface{})
		for sn, sv := range sub {
			tn := strings.ToUpper(sn)
			if tn == "OR" ||
				tn == "AND" {
				if sm, is := sv.(map[string]interface{}); is {
					subWhere[tn] = convertWhere(sm)
				}
				continue
			}
			if tn == "BETWEEN" {
				if vs, is := sv.([]interface{}); is {
					if len(vs) != 2 {
						continue
					}
					subWhere[tn] = vs
				}
				if vs, is := sv.([]string); is {
					var ds []interface{}
					for _, v := range vs {
						ds = append(ds, v)
					}
					if len(ds) != 2 {
						continue
					}
					subWhere[tn] = ds
				}
				continue
			}
			if tn == "IN" {
				if vs, is := sv.([]interface{}); is {
					if len(vs) == 0 {
						continue
					}
					subWhere[tn] = vs
				}
				if vs, is := sv.([]string); is {
					var ds []interface{}
					for _, v := range vs {
						ds = append(ds, v)
					}
					if len(ds) == 0 {
						continue
					}
					subWhere[tn] = ds
				}
				continue
			}
			subWhere[sn] = sv
		}
		where[utils.SnackedName(n)] = subWhere
	}
	return where
}

type DataTableSourceResponse struct {
	Total   int64       `json:"total"`
	PerPage int         `json:"perPage"`
	Page    int         `json:"page"`
	Data    interface{} `json:"data"`
}

// ParseTag parse tags setting
func (t *DataTableTzComponent) ParseTag(fields []*tzui.Field, tm *tzui.TagManager) {
	for i, field := range fields {
		to := TableOption{
			Field: field.CasedName(),
		}
		header, ok := field.Tags["header"]
		if !ok {
			// 必须项
			continue
		}
		to.Header = header

		if fieldType, ok := field.Tags["fieldType"]; ok {
			split := strings.Split(fieldType, ":")
			to.FieldType = split[0]
			if len(split) > 1 {
				if split[0] == "enum" || split[0] == "dic" {
					to.EnumName = split[1]
					if tm != nil {
						tm.GetTag("dictionary").IsValid(to.EnumName)
					}
				} else if split[0] == "text" {
					to.Placeholder = split[1]
				}
			}
		} else {
			to.FieldType = "text"
		}

		// 共享长度设置
		var twc *TableWidthConfig
		if width, ok := field.Tags["width"]; ok {
			twc = &TableWidthConfig{
				Field: field.CasedName(),
				Width: width,
			}
			t.TableWidthConfig = append(t.TableWidthConfig, *twc)
		} else {
			twc = &TableWidthConfig{
				Field:    field.CasedName(),
				Width:    "80px",
				widthNum: 80,
			}
			if to.FieldType == "date" ||
				to.FieldType == "enum" ||
				to.FieldType == "dic" {
				twc.Width = "120px"
				twc.widthNum = 120
			} else {
				// text
				twc.widthNum = 12*len([]rune(to.Header)) + 40
				twc.Width = strconv.Itoa(twc.widthNum) + "px"
			}
			t.TableWidthConfig = append(t.TableWidthConfig, *twc)
		}
		if fixedLeft, ok := field.Tags["fixedLeft"]; ok {
			if len(fixedLeft) == 0 {
				totalWidth := 0
				if i > 0 {
					for _, pre := range t.TableWidthConfig[:i] {
						totalWidth += pre.widthNum
					}
				}

				to.FixedLeft = strconv.Itoa(totalWidth) + "px"
			} else {
				to.FixedLeft = fixedLeft
			}
		}
		if resizeEnabled, ok := field.Tags["resizeEnabled"]; ok {
			if resizeEnabled == "false" {
				to.ResizeEnabled = false
			} else {
				to.ResizeEnabled = true
			}
		} else {
			to.ResizeEnabled = true
		}
		if form, ok := field.Tags["form"]; ok && len(form) > 0 {
			to.Form = form
		}
		if sort, ok := field.Tags["sort"]; ok {
			to.Sortable = true
			up := strings.ToUpper(sort)
			if up == "ASC" || up == "DESC" {
				to.Sort = up
			}
		}
		if percent, ok := field.Tags["percent"]; ok {
			to.Percent = percent
		}

		t.TableOptions = append(t.TableOptions, to)
	}
}

type dataTableTzComponentBuilder struct {
	model   interface{}
	typ     func() *DataTableTzComponent
	sources []*tzui.TzSource
}

func NewDataTableBuilder(model interface{}, source func(ctx context.Context, req *DataTableSourceRequest) (res *DataTableSourceResponse, err error), build func() *DataTableTzComponent) tzui.ITzComponentBuilder {
	if source == nil {
		panic("source can not be empty")
	}
	modelPkgPath := utils.GetPkgPath(model)
	fetchSource := &tzui.TzSource{
		Name:   "fetch",
		URL:    "fetch/" + modelPkgPath,
		Binder: (&DataTableSourceRequest{}).Bind,
		Handler: func(ctx context.Context, req interface{}) (interface{}, error) {
			return source(ctx, req.(*DataTableSourceRequest))
		},
	}

	return &dataTableTzComponentBuilder{
		model:   model,
		sources: []*tzui.TzSource{fetchSource},
		typ:     build,
	}
}

func (b dataTableTzComponentBuilder) Build() tzui.ITzComponent {
	c := b.typ()
	c.Sources = b.Sources()
	return c
}

func (b dataTableTzComponentBuilder) Model() interface{} {
	return b.model
}

func (b dataTableTzComponentBuilder) ComponentName() string {
	return dataTableTzComponentName
}

func (b dataTableTzComponentBuilder) Sources() []*tzui.TzSource {
	return b.sources
}

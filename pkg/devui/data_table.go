package devui

import (
	"context"
	"fmt"
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
	Scrollable             bool               `json:"scrollable"`
	Type                   string             `json:"type"`
	TableHeight            string             `json:"tableHeight"`
	VirtualScroll          bool               `json:"virtualScroll"`
	FixHeader              bool               `json:"fixHeader"`
	ContainFixHeaderHeight bool               `json:"containFixHeaderHeight"`
	TableWidthConfig       []TableWidthConfig `json:"tableWidthConfig"`
	TableOptions           []TableOption      `json:"tableOptions"`
}

type TableWidthConfig struct {
	Field string `json:"field"`
	Width string `json:"width"`
}

type TableOption struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	// typ 'text' | 'enum' 'dic' enum不显示value,dic显示[id]
	FieldType string `json:"fieldType"`
	EnumName  string
	// exp 120px
	FixedLeft     string `json:"fixedLeft"`
	ResizeEnabled bool   `json:"resizeEnable"`
	Sortable      bool
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

func (req *DataTableSourceRequest) Bind(body []byte) (interface{}, error) {
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
	var convertWhere func(src map[string]interface{}) map[string]interface{}
	convertWhere = func(src map[string]interface{}) map[string]interface{} {
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
				subWhere[utils.SnackedName(sn)] = sv
			}
			where[utils.SnackedName(n)] = v
		}
		return where
	}

	req.Where = convertWhere(req.Where)

	return req, err
}

type DataTableSourceResponse struct {
	Total   int64       `json:"total"`
	PerPage int         `json:"perPage"`
	Page    int         `json:"page"`
	Data    interface{} `json:"data"`
}

// ParseTag parse tags setting
func (t *DataTableTzComponent) ParseTag(fields []*tzui.Field, tm *tzui.TagManager) {
	for _, field := range fields {
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
			if len(split) == 1 {

				continue
			}
			to.FieldType = split[0]
			if split[0] == "enum" || split[0] == "dic" {
				to.EnumName = split[1]
				tm.GetTag("dictionary").IsValid(to.EnumName)
			}
		} else {
			to.FieldType = "text"
		}
		if fixedLeft, ok := field.Tags["fixedLeft"]; ok {
			to.FixedLeft = fixedLeft
		}
		if resizeEnabled, ok := field.Tags["resizeEnabled"]; ok {
			if resizeEnabled == "true" {
				to.ResizeEnabled = true
			}
		}
		if _, ok := field.Tags["sortable"]; ok {
			to.Sortable = true
		}
		if width, ok := field.Tags["width"]; ok {
			t.TableWidthConfig = append(t.TableWidthConfig, TableWidthConfig{
				Field: field.CasedName(),
				Width: width,
			})
		} else {
			t.TableWidthConfig = append(t.TableWidthConfig, TableWidthConfig{
				Field: field.CasedName(),
				Width: "80px",
			})
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

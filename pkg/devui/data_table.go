package devui

import (
	"context"
	"fmt"
	"reflect"

	"gitlab.com/tz/tzui"
)

var (
	_ tzui.HasSourceComponent        = &DataTableTzComponent{}
	_ tzui.HasSourceComponentBuilder = &dataTableTzComponentBuilder{}
	_ tzui.ITzComponentBuilder       = &dataTableTzComponentBuilder{}
)

// tzui tags:
// - header 表头名
// - fieldType 单元格格式 text
// - fieldLeft 左悬浮长度，用于固定列 0px 则第一列固定
// - width 单元格固定长度
type DataTableTzComponent struct {
	tzui.TzComponent
	DataSourceURL          string             `json:"dataSourceUrl"`
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
	// typ 'text' | ''
	FieldType string `json:"fieldType"`
	// exp 120px
	FixedLeft     string `json:"fixedLeft"`
	ResizeEnabled bool   `json:"resizeEnable"`
}

type DataTableSourceRequest struct {
	Where   map[string]interface{} `json:"where"`
	Total   int64                  `json:"total"`
	PerPage int                    `json:"perPage"`
	Page    int                    `json:"page"`
}

func (req *DataTableSourceRequest) Data(data interface{}) *DataTableSourceResponse {
	return &DataTableSourceResponse{
		Total:   req.Total,
		PerPage: req.PerPage,
		Page:    req.Page,
		Data:    data,
	}
}

type DataTableSourceResponse struct {
	Total   int64       `json:"total"`
	PerPage int         `json:"perPage"`
	Page    int         `json:"page"`
	Data    interface{} `json:"data"`
}

// ReqStr return request struct with default value
func (t DataTableTzComponent) ReqStr() interface{} {
	return &DataTableSourceRequest{
		Total: 50,
	}
}

// ResStr return response struct with default value
func (t DataTableTzComponent) ResStr() interface{} {
	return &DataTableSourceResponse{Total: 50}
}

func (t *DataTableTzComponent) SetSource(source string) {
	t.DataSourceURL = source
}

// ParseTag parse tags setting
func (t *DataTableTzComponent) ParseTag(fields []*tzui.Field) {
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
			to.FieldType = fieldType
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

	}
}

type dataTableTzComponentBuilder struct {
	model  interface{}
	source func(ctx context.Context, req *DataTableSourceRequest) (res *DataTableSourceResponse, err error)
	typ    func() *DataTableTzComponent
}

func NewDataTableBuilder(model interface{}, source func(ctx context.Context, req *DataTableSourceRequest) (res *DataTableSourceResponse, err error), build func() *DataTableTzComponent) tzui.ITzComponentBuilder {
	if source == nil {
		panic("source can not be empty")
	}
	return &dataTableTzComponentBuilder{
		model:  model,
		source: source,
		typ:    build,
	}
}

func (b dataTableTzComponentBuilder) Build() tzui.ITzComponent {
	return b.typ()
}

func (b dataTableTzComponentBuilder) Model() interface{} {
	return b.model
}

func (b dataTableTzComponentBuilder) ComponentName() string {
	return dataTableTzComponentName
}

func (b *dataTableTzComponentBuilder) Handler(ctx context.Context, req interface{}) (interface{}, error) {
	return b.source(ctx, req.(*DataTableSourceRequest))
}

func (b dataTableTzComponentBuilder) SourceURL() string {
	sourceType := reflect.TypeOf(b.source)
	return fmt.Sprintf("%s/%s", sourceType.PkgPath(), sourceType.Name())
}

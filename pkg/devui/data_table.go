package devui

import "gitlab.com/tz/tzui/pkg/tzui"

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
	Where   map[string]string `json:"where"`
	Total   int64             `json:"total"`
	PerPage int               `json:"perPage"`
	Page    int               `json:"page"`
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
